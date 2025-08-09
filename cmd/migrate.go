package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/turahe/master-data-rest-api/internal/adapters/secondary/database"
)

var (
	migrationsDir    string
	migrationStep    int
	migrationVersion int
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration management",
	Long: `Manage database migrations for the Master Data REST API.

This command provides various migration operations including:
- Running all pending migrations (up)
- Rolling back migrations (down)
- Checking migration status
- Forcing to a specific migration version

Examples:
  # Run all pending migrations
  master-data-api migrate up

  # Rollback the last migration
  master-data-api migrate down

  # Check migration status
  master-data-api migrate status

  # Rollback specific number of migrations
  master-data-api migrate down --step 2`,
}

// migrateUpCmd represents the migrate up command
var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Run pending database migrations",
	Long: `Run all pending database migrations to bring the database
schema up to the latest version.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := CheckConfig(cmd); err != nil {
			return err
		}
		return runMigrationsUp()
	},
}

// migrateDownCmd represents the migrate down command
var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback database migrations",
	Long: `Rollback database migrations. By default, rolls back one migration.
Use --step flag to specify the number of migrations to rollback.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := CheckConfig(cmd); err != nil {
			return err
		}
		return runMigrationsDown()
	},
}

// migrateStatusCmd represents the migrate status command
var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	Long:  `Display the current status of database migrations.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := CheckConfig(cmd); err != nil {
			return err
		}
		return showMigrationStatus()
	},
}

// migrateForceCmd represents the migrate force command
var migrateForceCmd = &cobra.Command{
	Use:   "force",
	Short: "Force migration to a specific version",
	Long: `Force the migration version to a specific number. This is useful when
migrations are in a dirty state and need to be manually fixed.

WARNING: Use this command with caution as it can lead to data loss if used incorrectly.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := CheckConfig(cmd); err != nil {
			return err
		}
		return forceMigrationVersion()
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateStatusCmd)
	migrateCmd.AddCommand(migrateForceCmd)

	// Migration flags
	migrateCmd.PersistentFlags().StringVarP(&migrationsDir, "migrations-dir", "m", "migrations", "directory containing migration files")
	migrateDownCmd.Flags().IntVarP(&migrationStep, "step", "s", 1, "number of migrations to rollback")
	migrateForceCmd.Flags().IntVarP(&migrationVersion, "version", "v", 1, "version to force migration to")
}

func runMigrationsUp() error {
	config := GetConfig()
	log := GetLogger()

	log.Info("Running database migrations up")

	migrator := database.NewMigrator(config.Database)
	if err := migrator.RunMigrations(migrationsDir); err != nil {
		log.WithError(err).Error("Failed to run migrations")
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Info("Database migrations completed successfully")
	fmt.Println("✅ All migrations completed successfully!")

	return nil
}

func runMigrationsDown() error {
	config := GetConfig()
	log := GetLogger()

	log.WithField("step", migrationStep).Info("Rolling back database migrations")

	migrator := database.NewMigrator(config.Database)
	if err := migrator.RollbackMigrations(migrationsDir, migrationStep); err != nil {
		log.WithError(err).Error("Failed to rollback migrations")
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	log.WithField("step", migrationStep).Info("Database migrations rolled back successfully")
	fmt.Printf("✅ Successfully rolled back %d migration(s)!\n", migrationStep)

	return nil
}

func showMigrationStatus() error {
	config := GetConfig()
	log := GetLogger()

	log.Info("Checking migration status")

	migrator := database.NewMigrator(config.Database)
	version, dirty, err := migrator.GetMigrationStatus()
	if err != nil {
		log.WithError(err).Error("Failed to get migration status")
		return fmt.Errorf("failed to get migration status: %w", err)
	}

	fmt.Println("Database Migration Status")
	fmt.Println("========================")
	fmt.Printf("Current Version: %d\n", version)
	fmt.Printf("Dirty State:     %v\n", dirty)

	if dirty {
		fmt.Println("⚠️  Database is in a dirty state. This usually means a migration failed.")
		fmt.Println("   You may need to manually fix the database or force to a specific version.")
	} else {
		fmt.Println("✅ Database is in a clean state.")
	}

	return nil
}

func forceMigrationVersion() error {
	config := GetConfig()
	log := GetLogger()

	log.WithField("version", migrationVersion).Warn("Forcing migration to specific version - USE WITH CAUTION!")

	migrator := database.NewMigrator(config.Database)
	if err := migrator.ForceMigrationVersion(migrationsDir, migrationVersion); err != nil {
		log.WithError(err).Error("Failed to force migration version")
		return fmt.Errorf("failed to force migration version: %w", err)
	}

	log.WithField("version", migrationVersion).Info("Migration version forced successfully")
	fmt.Printf("✅ Successfully forced migration to version %d!\n", migrationVersion)
	fmt.Println("⚠️  Remember to run 'migrate status' to verify the state.")

	return nil
}
