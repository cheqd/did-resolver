package migrations

func MigrateDID(did string) string {
	did = MigrateIndyStyleDid(did)
	did = MigrateUUIDDid(did)
	return did
}
