{
  description = "context - devops cli project called to solve the problem of fucking with legacy code";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = [
            pkgs.go
            pkgs.gnumake
            pkgs.git
            pkgs.gopls
            pkgs.gofumpt
          ];

          shellHook = ''
            echo "context hash   : (generated at build time by internal/auth)"
            echo "go version     : $(go version)"
            echo ""
            echo "make targets   :"
            echo "make build     : build for nixos"
            echo "make build-win : cross-compile for windows"
            echo "make run       : build and run on current project"
          '';
        };

        packages.default = pkgs.buildGoModule {
          pname = "context";
          version = "0.0";
          src = ./.;
          vendorHash = null;
        };
      }
    );
}
