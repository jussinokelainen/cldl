class Cldl < Formula
  desc "Directory based cli todo lists"
  homepage "https://github.com/jussinokelainen/cldl"
  url "https://github.com/jussinokelainen/cldl/archive/refs/tags/v0.2.2.tar.gz"
  sha256 "dddc6093ac82d7193e0a239da7fbea63e9acb2dcd568e692c4e6287fd87ccf70"
  license "MIT"

  depends_on "go" => :build

  def install
    system "make", "release"
    bin.install "build/cldl"
  end
end
