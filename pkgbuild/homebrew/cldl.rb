class Cldl < Formula
  desc "Directory based cli todo lists"
  homepage "https://github.com/jussinokelainen/cldl"
  url "https://github.com/jussinokelainen/cldl/archive/refs/tags/v0.2.3.tar.gz"
  sha256 "7cb3d3fddbfae8f24e4b0340eec0373860d8cbd92df516d8530d6f5c38351bf9"
  license "MIT"

  depends_on "go" => :build

  def install
    system "make", "release"
    bin.install "build/cldl"
  end
end
