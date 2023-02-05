load("@bazel_gazelle//:def.bzl", "DEFAULT_LANGUAGES", "gazelle", "gazelle_binary")

# gazelle:go_proto_compilers @io_bazel_rules_go//proto:go_proto
# gazelle:go_grpc_compilers //:gen-go-grpc
# gazelle:prefix github.com/kchymet/generic-grpc
gazelle(
    name = "gazelle",
)

load("@io_bazel_rules_go//proto/wkt:well_known_types.bzl", "WELL_KNOWN_TYPES_APIV2")
load("@io_bazel_rules_go//proto:compiler.bzl", "go_proto_compiler")

# Borrowed rule from sourcegraph PR:
# Because the current implementation of rules_go uses the old protoc grpc compiler, we have to declare our own, and declare it manually in the build files.
# See https://github.com/bazelbuild/rules_go/issues/3022
go_proto_compiler(
    name = "gen-go-grpc",
    plugin = "@org_golang_google_grpc_cmd_protoc_gen_go_grpc//:protoc-gen-go-grpc",
    suffix = "_grpc.pb.go",
    valid_archive = False,
    visibility = ["//visibility:public"],
    deps = WELL_KNOWN_TYPES_APIV2 + [
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)
