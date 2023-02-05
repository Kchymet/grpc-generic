load("@bazel_gazelle//:def.bzl", "DEFAULT_LANGUAGES", "gazelle", "gazelle_binary")

gazelle_binary(
    name = "gazelle_binary",
    languages = DEFAULT_LANGUAGES,
    visibility = ["//visibility:private"],
)

# gazelle:prefix github.com/kchymet/generic-grpc
gazelle(
    name = "gazelle",
    gazelle = "//:gazelle_binary",
)
