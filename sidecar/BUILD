load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["sidecar.go"],
    importpath = "github.com/buildbuddy-io/buildbuddy-cli/sidecar",
    visibility = ["//visibility:public"],
    deps = [
        "//download:go_default_library",
        "//logging:go_default_library",
        "@org_golang_x_mod//semver:go_default_library",
    ],
)
