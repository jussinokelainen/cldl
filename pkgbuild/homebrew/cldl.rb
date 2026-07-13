class Cldl < Formula
  desc "Todo CLI"
  homepage "https://github.com/jussinokelainen/cldl"
  url "https://github.com/jussinokelainen/cldl/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "8e7c510772da3c2efe889c8de1926a1436b71152c89ed60ce50879d4ca2c7c44"
  license "MIT"

  depends_on "go" => :build

  def install
    system "make", "release"
    bin.install "bin/cldl"
  end
end
