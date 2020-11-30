load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["sidecar.go"],
    importpath = "github.com/buildbuddy-io/buildbuddy-cli/cmd/sidecar",
    visibility = ["//visibility:private"],
    deps = [
        "//cache_proxy:go_default_library",
        "//devnull:go_default_library",
        "@com_github_buildbuddy_io_buildbuddy//proto:publish_build_event_go_proto",
        "@com_github_buildbuddy_io_buildbuddy//proto:remote_execution_go_proto",
        "@com_github_buildbuddy_io_buildbuddy//server/backends/disk_cache:go_default_library",
        "@com_github_buildbuddy_io_buildbuddy//server/build_event_protocol/build_event_proxy:go_default_library",
        "@com_github_buildbuddy_io_buildbuddy//server/build_event_protocol/build_event_server:go_default_library",
        "@com_github_buildbuddy_io_buildbuddy//server/config:go_default_library",
        "@com_github_buildbuddy_io_buildbuddy//server/nullauth:go_default_library",
        "@com_github_buildbuddy_io_buildbuddy//server/real_environment:go_default_library",
        "@com_github_buildbuddy_io_buildbuddy//server/rpc/filters:go_default_library",
        "@com_github_buildbuddy_io_buildbuddy//server/util/grpc_client:go_default_library",
        "@com_github_buildbuddy_io_buildbuddy//server/util/grpc_server:go_default_library",
        "@com_github_buildbuddy_io_buildbuddy//server/util/healthcheck:go_default_library",
        "@go_googleapis//google/bytestream:bytestream_go_proto",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//reflection:go_default_library",
    ],
)

go_binary(
    name = "sidecar",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

go_binary(
    name = "sidecar-darwin",
    out = "sidecar-darwin-amd64",
    embed = [":go_default_library"],
    gc_linkopts = [
        "-s",
        "-w",
    ],
    goarch = "amd64",
    goos = "darwin",
    pure = "on",
    tags = ["manual"],
    visibility = ["//visibility:public"],
)

go_binary(
    name = "sidecar-linux",
    out = "sidecar-linux-amd64",
    embed = [":go_default_library"],
    gc_linkopts = [
        "-s",
        "-w",
    ],
    goarch = "amd64",
    goos = "linux",
    pure = "on",
    tags = ["manual"],
    visibility = ["//visibility:public"],
)

go_binary(
    name = "sidecar-linux-arm64",
    out = "sidecar-linux-arm64",
    embed = [":go_default_library"],
    gc_linkopts = [
        "-s",
        "-w",
    ],
    goarch = "arm64",
    goos = "linux",
    pure = "on",
    tags = ["manual"],
    visibility = ["//visibility:public"],
)

go_binary(
    name = "sidecar-windows",
    out = "sidecar-windows-amd64.exe",
    embed = [":go_default_library"],
    goarch = "amd64",
    goos = "windows",
    pure = "on",
    tags = ["manual"],
    visibility = ["//visibility:public"],
)
