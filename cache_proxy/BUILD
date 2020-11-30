load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["cache_proxy.go"],
    importpath = "github.com/buildbuddy-io/buildbuddy-cli/cache_proxy",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_buildbuddy_io_buildbuddy//proto:remote_execution_go_proto",
        "@com_github_buildbuddy_io_buildbuddy//server/environment:go_default_library",
        "@com_github_buildbuddy_io_buildbuddy//server/interfaces:go_default_library",
        "@com_github_buildbuddy_io_buildbuddy//server/remote_cache/byte_stream_server:go_default_library",
        "@com_github_buildbuddy_io_buildbuddy//server/remote_cache/cachetools:go_default_library",
        "@com_github_buildbuddy_io_buildbuddy//server/remote_cache/content_addressable_storage_server:go_default_library",
        "@com_github_buildbuddy_io_buildbuddy//server/remote_cache/digest:go_default_library",
        "@com_github_buildbuddy_io_buildbuddy//server/util/prefix:go_default_library",
        "@com_github_buildbuddy_io_buildbuddy//server/util/status:go_default_library",
        "@go_googleapis//google/bytestream:bytestream_go_proto",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//test/bufconn:go_default_library",
    ],
)
