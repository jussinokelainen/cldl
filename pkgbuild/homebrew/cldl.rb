class Cldl < Formula
  desc "Todo CLI"
  homepage "https://github.com/jussinokelainen/cldl"
  url "https://github.com/jussinokelainen/cldl/archive/refs/tags/v0.2.0.tar.gz"
  sha256 "855d19ca0812b06936cf972c126bebb9f099ce030ac6bee7538198a7e95a7a13"
  license "MIT"

  depends_on "go" => :build

  def install
    system "make", "release"
    bin.install "bin/cldl"
  end
end
