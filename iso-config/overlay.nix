self: super: {
  ztrace = super.stdenv.mkDerivation {
    pname = "ztrace";
    version = "1.0";

    src = /home/haiba/ztrace; # must be an absolute path UPDATEME

    dontUnpack = true;

    installPhase = ''
      mkdir -p $out/bin
      cp $src $out/bin/ztrace
      chmod +x $out/bin/ztrace
    '';
  };
}
