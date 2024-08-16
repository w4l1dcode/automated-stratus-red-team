package stratus

import (
	"database/sql"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/datadog/stratus-red-team/v2/pkg/stratus"
	_ "github.com/datadog/stratus-red-team/v2/pkg/stratus/mitreattack"
	stratusrunner "github.com/datadog/stratus-red-team/v2/pkg/stratus/runner"
	"github.com/sirupsen/logrus"
)

func DetonateTTPs(db *sql.DB, platform string, awsClient *ecr.Client, l *logrus.Logger) error {
	tactic := GetUnusedTactic(db)

	filter := &stratus.AttackTechniqueFilter{
		Platform: stratus.Platform(platform),
		Tactic:   tactic,
	}

	ttps := stratus.GetRegistry().GetAttackTechniques(filter)

	logrus.Infof("Number of ttps found for %s for tactic \"%s\": %d \n", platform, TacticToString(tactic), len(ttps))

	if len(ttps) == 0 {
		logrus.Warningf("No TTPs found for tactic: %s\n", TacticToString(tactic))
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
				logrus.Warningf("Could not cleanup created infrastructure: %v\n", err)
			}
		}(stratusRunner)

		logrus.Infof("TTP %s is warm! Executing...\n", ttp.ID)

		err = stratusRunner.Detonate()

		if err != nil {
			logrus.Warningf("Could not detonate TTP: %v\n", err)
			logrus.Warning("Continuing...")
		} else {
			logrus.Infof("TTP %s detonated\n", ttp.ID)
		}

		// Add a 20-minute delay between attacks to wait if a detection will be triggered
		//time.Sleep(20 * time.Minute)
	}

	err := MarkTacticAsUsed(db, TacticToString(tactic))
	if err != nil {
		logrus.Fatalf("Failed marking tactic as used: %s\n", err)
	}
	return nil
}
