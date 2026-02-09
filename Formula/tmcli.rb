# Formula for tmcli - Time Machine CLI
# https://github.com/nvanwyen/tmcli
#
# To install:
#   brew tap nvanwyen/tmcli
#   brew install tmcli

class Tmcli < Formula
  desc "macOS Time Machine CLI and interactive TUI wrapping tmutil"
  homepage "https://github.com/nvanwyen/tmcli"
  url "https://github.com/nvanwyen/tmcli/archive/refs/tags/v1.0.2.tar.gz"
  sha256 "REPLACE_WITH_ACTUAL_SHA256"
  license "MIT"

  depends_on "cmake" => :build
  depends_on "go" => :build
  depends_on :macos

  def install
    # Configure and build using the project's CMake build system
    mkdir "build" do
      system "cmake", "..", "-DCMAKE_BUILD_TYPE=Release"
      system "make"
    end

    # Install the binary
    bin.install "bin/tmcli"
  end

  test do
    assert_match "Time Machine CLI", shell_output("#{bin}/tmcli --version")
  end
end
