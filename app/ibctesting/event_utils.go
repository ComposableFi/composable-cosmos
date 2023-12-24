package ibctesting

import (
	"fmt"
	"strconv"
	"strings"

	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
)

func getSendPackets(evts []abci.Event) []channeltypes.Packet {
	var res []channeltypes.Packet
	for _, evt := range evts {
		if evt.Type == "send_packet" {
			packet := parsePacketFromEvent(evt)
			res = append(res, packet)
		}
	}
	return res
}

func getAckPackets(evts []abci.Event) []PacketAck {
	var res []PacketAck
	for _, evt := range evts {
		if evt.Type == "write_acknowledgement" {
			packet := parsePacketFromEvent(evt)
			ack := PacketAck{
				Packet: packet,
				Ack:    []byte(getField(evt, "packet_ack")),
			}
			res = append(res, ack)
		}
	}
	return res
}

// Used for various debug statements above when needed... do not remove
// func showEvent(evt abci.Event) {
//	fmt.Printf("evt.Type: %s\n", evt.Type)
//	for _, attr := range evt.Attributes {
//		fmt.Printf("  %s = %s\n", string(attr.Key), string(attr.Value))
//	}
//}

func parsePacketFromEvent(evt abci.Event) channeltypes.Packet {
	return channeltypes.Packet{
		Sequence:           getUintField(evt, "packet_sequence"),
		SourcePort:         getField(evt, "packet_src_port"),
		SourceChannel:      getField(evt, "packet_src_channel"),
		DestinationPort:    getField(evt, "packet_dst_port"),
		DestinationChannel: getField(evt, "packet_dst_channel"),
		Data:               []byte(getField(evt, "packet_data")),
		TimeoutHeight:      parseTimeoutHeight(getField(evt, "packet_timeout_height")),
		TimeoutTimestamp:   getUintField(evt, "packet_timeout_timestamp"),
	}
}

// ParsePacketFromEvents parses events emitted from a MsgRecvPacket and returns the
// acknowledgement.
func ParsePacketFromEvents(events sdk.Events) (channeltypes.Packet, error) {
	for _, ev := range events {
		if ev.Type == channeltypes.EventTypeSendPacket {
			packet := channeltypes.Packet{}
			for _, attr := range ev.Attributes {
				switch attr.Key {
				case channeltypes.AttributeKeyData: //nolint: staticcheck
					packet.Data = []byte(attr.Value)

				case channeltypes.AttributeKeySequence:
					seq, err := strconv.ParseUint(attr.Value, 10, 64)
					if err != nil {
						return channeltypes.Packet{}, err
					}

					packet.Sequence = seq

				case channeltypes.AttributeKeySrcPort:
					packet.SourcePort = attr.Value

				case channeltypes.AttributeKeySrcChannel:
					packet.SourceChannel = attr.Value

				case channeltypes.AttributeKeyDstPort:
					packet.DestinationPort = attr.Value

				case channeltypes.AttributeKeyDstChannel:
					packet.DestinationChannel = attr.Value

				case channeltypes.AttributeKeyTimeoutHeight:
					height, err := clienttypes.ParseHeight(attr.Value)
					if err != nil {
						return channeltypes.Packet{}, err
					}

					packet.TimeoutHeight = height

				case channeltypes.AttributeKeyTimeoutTimestamp:
					timestamp, err := strconv.ParseUint(attr.Value, 10, 64)
					if err != nil {
						return channeltypes.Packet{}, err
					}

					packet.TimeoutTimestamp = timestamp

				default:
					continue
				}
			}

			return packet, nil
		}
	}
	return channeltypes.Packet{}, fmt.Errorf("acknowledgement event attribute not found")
}

// return the value for the attribute with the given name
func getField(evt abci.Event, key string) string {
	for _, attr := range evt.Attributes {
		if attr.Key == key {
			return attr.Value
		}
	}
	return ""
}

func getUintField(evt abci.Event, key string) uint64 {
	raw := getField(evt, key)
	return toUint64(raw)
}

func toUint64(raw string) uint64 {
	if raw == "" {
		return 0
	}
	i, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func parseTimeoutHeight(raw string) clienttypes.Height {
	chunks := strings.Split(raw, "-")
	return clienttypes.Height{
		RevisionNumber: toUint64(chunks[0]),
		RevisionHeight: toUint64(chunks[1]),
	}
}
