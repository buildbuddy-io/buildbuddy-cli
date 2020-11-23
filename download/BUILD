load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["download.go"],
    importpath = "github.com/buildbuddy-io/buildbuddy-cli/download",
    visibility = ["//visibility:public"],
    deps = [
        "//logging:go_default_library",
        "@com_github_buildbuddy_io_buildbuddy//server/util/status:go_default_library",
        "@com_github_google_go_github//github:go_default_library",
    ],
)
