{
  description = "Dev Flake !";

  inputs.nixpkgs.url = "https://flakehub.com/f/NixOS/nixpkgs/0.1.*.tar.gz";

  outputs = { self, nixpkgs }:
    let
      # goVersion = 22; # Change this to update the whole stack

      supportedSystems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forEachSupportedSystem = f: nixpkgs.lib.genAttrs supportedSystems (system: f {
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ self.overlays.default ];
        };
      });
    in
    {
      overlays.default = final: prev: {
        # go = final."go_1_${toString goVersion}";
      };

      devShells = forEachSupportedSystem ({ pkgs }: {
        default = pkgs.mkShell {
          # packages = with pkgs; [ # things needet at coding time
          # ];
          # buildInputs = with pkgs; [ #things needed at runtime
          # ];
          nativeBuildInputs = with pkgs; [ #things needed at compile time
            #pkg-config
            # alsa-lib
            #SDL2
          ];
        };
      });
    };
}
