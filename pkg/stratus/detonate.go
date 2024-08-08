package stratus

import (
	"database/sql"
	"github.com/datadog/stratus-red-team/v2/pkg/stratus"
	_ "github.com/datadog/stratus-red-team/v2/pkg/stratus/mitreattack"
	stratusrunner "github.com/datadog/stratus-red-team/v2/pkg/stratus/runner"
	"github.com/sirupsen/logrus"
	"time"
)

func DetonateTTPs(db *sql.DB, platform string) error {
	tactic := GetUnusedTactic(db)

	logrus.Infof("Executing tests of tactic: %s\n", tactic)

	filter := &stratus.AttackTechniqueFilter{
		Platform: stratus.Platform(platform),
		Tactic:   tactic,
	}

	ttps := stratus.GetRegistry().GetAttackTechniques(filter)

	logrus.Infof("Number of ttps found: %d\n", len(ttps))

	if len(ttps) == 0 {
		logrus.Fatalf("No ttps found for tactic: %s\n", TacticToString(tactic))
		return nil
	}

	for _, ttp := range ttps {
		logrus.Infof("Executing ttp: %s\n", ttp.ID)

		stratusRunner := stratusrunner.NewRunner(ttp, stratusrunner.StratusRunnerNoForce)
		if stratusRunner == nil {
			logrus.Fatalf("Failed to create StratusRunner for ttp: %s\n", ttp.ID)
		}

		_, err := stratusRunner.WarmUp()
		if err != nil {
			logrus.Fatalf("Could not warm up TTP %s\n: %v", ttp.ID, err)
		}
		defer func(stratusRunner stratusrunner.Runner) {
			err := stratusRunner.CleanUp()
			if err != nil {
				logrus.Infof("Could not cleanup created infrastructure: %v\n", err)
			}
		}(stratusRunner)

		logrus.Info("TTP is warm! Executing...\n")

		err = stratusRunner.Detonate()
		if err != nil {
			logrus.Fatalf("Could not detonate TTP: %s\n", err)
		}

		logrus.Info("TTP detonated!\n")

		// Add a 20-minute delay between attacks to wait if a detection will be triggered
		time.Sleep(20 * time.Minute)
	}

	MarkTacticAsUsed(db, TacticToString(tactic))
	return nil
}
