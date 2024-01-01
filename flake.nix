{
  description = "composable-cosmos";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    gomod2nix = {
      url = github:nix-community/gomod2nix;
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };
  outputs = inputs @ {
    flake-parts,
    gomod2nix,
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
      }: {
        formatter = pkgs.alejandra;
        devShells = {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [
              pkgs.go
              pkgs.gnumake
              pkgs.gotools
              pkgs.golangci-lint
              pkgs.gci
              gomod2nix.packages.${system}.default
            ];
            shellHook = ''
            '';
          };
        };
        packages.default = pkgs.writeShellApplication {
          text = ''
              go get mvdan.cc/gofumpt
              go get github.com/client9/misspell/cmd/misspell
              go get golang.org/x/tools/cmd/goimports
              make lint
              nix fmt
              go mod tidy
              golangci-lint run ./... --fix
          '';
        };
      };
    };
}
