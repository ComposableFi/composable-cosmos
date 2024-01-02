{
  description = "composable-cosmos";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    gomod2nix = {
      url = github:nix-community/gomod2nix;
      inputs.nixpkgs.follows = "nixpkgs";
    };
    cosmos = {
      url = "github:informalsystems/cosmos.nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };
  outputs = inputs @ {
    flake-parts,
    gomod2nix,
    cosmos,
    ...
  }:
    flake-parts.lib.mkFlake {inherit inputs;} {
      imports = [
      ];
      systems = ["x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin"];
      perSystem = {
        config,
        self',
        inputs',
        pkgs,
        system,
        ...
      }: let
        inherit (inputs.gomod2nix.legacyPackages.${system}) buildGoApplication mkGoEvn gomod2nix;
        libwasmvm = cosmos.packages.${system}.libwasmvm_1_5_0;
        cosmwasm-check = cosmos.packages.${system}.cosmwasm-check;
        cosmosLib = cosmos.lib {inherit pkgs cosmwasm-check;};
      in {
        formatter = pkgs.alejandra;
        devShells = {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [
              bingo
              gci
              gnumake
              go
              golangci-lint
              gomod2nix
              gotools
            ];
          };
        };

        packages = {
          centaurid = buildGoApplication {
            pname = "centaurid";
            version = "v7.0.0";
            src = ./.;
            subPackages = ["cmd/centaurid"];
            modules = ./gomod2nix.toml;
            preFixup = ''
              ${cosmosLib.wasmdPreFixupPhase libwasmvm "centaurid"}
            '';
            buildInputs = [libwasmvm];
          };
          default = pkgs.writeShellApplication {
            name = "ci";
            text = ''
              go get mvdan.cc/gofumpt
              go get github.com/client9/misspell/cmd/misspell
              go get golang.org/x/tools/cmd/goimports
              make lint
              nix fmt
              go mod tidy
              golangci-lint run --fix
              make build
              gomod2nix
            '';
          };
        };
      };
    };
}
