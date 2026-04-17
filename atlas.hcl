env "local" {
  src = "ent://internal/repo/persistent/ent/schema"
  dev = "docker://postgres/16/dev?search_path=public"

  migration {
    dir = "file://migrations"
  }
}
