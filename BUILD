package(default_visibility = ["//visibility:public"])

load("@bazel_gazelle//:def.bzl", "gazelle")

# Ignore the node_modules dir
# gazelle:exclude node_modules
# Prefer generated BUILD files to be called BUILD over BUILD.bazel
# gazelle:build_file_name BUILD,BUILD.bazel
# gazelle:prefix github.com/buildbuddy-io/buildbuddy-cli
# gazelle:resolve go github.com/buildbuddy-io/buildbuddy/proto/invocation @com_github_buildbuddy_io_buildbuddy//proto:invocation_go_proto
# gazelle:resolve go github.com/buildbuddy-io/buildbuddy/proto/remote_execution @com_github_buildbuddy_io_buildbuddy//proto:remote_execution_go_proto
# gazelle:resolve go github.com/buildbuddy-io/buildbuddy/proto/build_event_stream @com_github_buildbuddy_io_buildbuddy//proto:build_event_stream_go_proto
# gazelle:resolve go github.com/buildbuddy-io/buildbuddy/proto/api/v1 @com_github_buildbuddy_io_buildbuddy//proto/api/v1:api_v1_go_proto
# gazelle:resolve go github.com/buildbuddy-io/buildbuddy/proto/telemetry @com_github_buildbuddy_io_buildbuddy//proto:telemetry_go_proto
# gazelle:resolve go github.com/buildbuddy-io/buildbuddy/proto/execution_stats @com_github_buildbuddy_io_buildbuddy//proto:execution_stats_go_proto
# gazelle:resolve go github.com/buildbuddy-io/buildbuddy/proto/scheduler @com_github_buildbuddy_io_buildbuddy//proto:scheduler_go_proto
# gazelle:resolve go github.com/buildbuddy-io/buildbuddy/proto/group @com_github_buildbuddy_io_buildbuddy//proto:group_go_proto
gazelle(name = "gazelle")

exports_files([
    "VERSION",
])
