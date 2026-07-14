class Cldl < Formula
  desc "Directory based cli todo lists"
  homepage "https://github.com/jussinokelainen/cldl"
  url "https://github.com/jussinokelainen/cldl/archive/refs/tags/v0.2.1.tar.gz"
  sha256 "ed41bc1aaf5d79c189db7a220266e4bf4d87f19d605d5ace0942bdafb037ed19"
  license "MIT"

  depends_on "go" => :build

  def install
    system "make", "release"
    bin.install "bin/cldl"
  end
end
