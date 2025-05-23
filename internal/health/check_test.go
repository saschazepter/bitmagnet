package health

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatusUnknownBeforeStatusUp(t *testing.T) {
	t.Parallel()

	// Arrange
	testData := map[string]CheckState{"check1": {Status: StatusUp}, "check2": {Status: StatusUnknown}}

	// Act
	result := aggregateStatus(testData)

	// Assert
	assert.Equal(t, StatusUnknown, result)
}

func TestStatusDownBeforeStatusUnknown(t *testing.T) {
	t.Parallel()

	// Arrange
	testData := map[string]CheckState{"check1": {Status: StatusDown}, "check2": {Status: StatusUnknown}}

	// Act
	result := aggregateStatus(testData)

	// Assert
	assert.Equal(t, StatusDown, result)
}

func doTestEvaluateAvailabilityStatus(
	t *testing.T,
	expectedStatus AvailabilityStatus,
	maxTimeInError time.Duration,
	maxFails uint,
	state CheckState,
) {
	t.Helper()

	// Act
	result := evaluateCheckStatus(&state, maxTimeInError, maxFails)

	// Assert
	assert.Equal(t, expectedStatus, result)
}

func TestWhenNoChecksMadeYetThenStatusUnknown(t *testing.T) {
	t.Parallel()

	doTestEvaluateAvailabilityStatus(t, StatusUnknown, 0, 0, CheckState{
		LastCheckedAt: time.Time{},
	})
}

func TestWhenNoErrorThenStatusUp(t *testing.T) {
	t.Parallel()

	doTestEvaluateAvailabilityStatus(t, StatusUp, 0, 0, CheckState{
		LastCheckedAt: time.Now(),
	})
}

func TestWhenErrorThenStatusDown(t *testing.T) {
	t.Parallel()

	doTestEvaluateAvailabilityStatus(t, StatusDown, 0, 0, CheckState{
		LastCheckedAt: time.Now(),
		Result:        fmt.Errorf("example error"),
	})
}

func TestWhenErrorAndMaxFailuresThresholdNotCrossedThenStatusWarn(t *testing.T) {
	t.Parallel()

	doTestEvaluateAvailabilityStatus(t, StatusUp, 1*time.Second, uint(10), CheckState{
		LastCheckedAt:       time.Now(),
		Result:              fmt.Errorf("example error"),
		FirstCheckStartedAt: time.Now().Add(-2 * time.Minute),
		LastSuccessAt:       time.Now().Add(-3 * time.Minute),
		ContiguousFails:     1,
	})
}

func TestWhenErrorAndMaxTimeInErrorThresholdNotCrossedThenStatusWarn(t *testing.T) {
	t.Parallel()

	doTestEvaluateAvailabilityStatus(t, StatusUp, 1*time.Hour, uint(1), CheckState{
		LastCheckedAt:       time.Now(),
		Result:              fmt.Errorf("example error"),
		FirstCheckStartedAt: time.Now().Add(-3 * time.Minute),
		LastSuccessAt:       time.Now().Add(-2 * time.Minute),
		ContiguousFails:     100,
	})
}

func TestWhenErrorAndAllThresholdsCrossedThenStatusDown(t *testing.T) {
	t.Parallel()

	doTestEvaluateAvailabilityStatus(t, StatusDown, 1*time.Second, uint(1), CheckState{
		LastCheckedAt:       time.Now(),
		Result:              fmt.Errorf("example error"),
		FirstCheckStartedAt: time.Now().Add(-3 * time.Minute),
		LastSuccessAt:       time.Now().Add(-2 * time.Minute),
		ContiguousFails:     5,
	})
}

func TestStartStopManualPeriodicChecks(t *testing.T) {
	t.Parallel()

	ckr := NewChecker(
		WithDisabledAutostart(),
		WithPeriodicCheck(50*time.Minute, 0, Check{
			Name: "check",
			Check: func(context.Context) error {
				return nil
			},
		}))

	assert.Equal(t, 0, ckr.GetRunningPeriodicCheckCount())

	ckr.Start()
	assert.Equal(t, 1, ckr.GetRunningPeriodicCheckCount())

	ckr.Stop()
	assert.Equal(t, 0, ckr.GetRunningPeriodicCheckCount())
}

func doTestCheckerCheckFunc(t *testing.T, updateInterval time.Duration, err error, expectedStatus AvailabilityStatus) {
	t.Helper()

	// Arrange
	ckr := NewChecker(
		WithTimeout(10*time.Second),
		WithCheck(Check{
			Name: "check1",
			Check: func(context.Context) error {
				return nil
			},
		}),
		WithPeriodicCheck(updateInterval, 0, Check{
			Name: "check2",
			Check: func(context.Context) error {
				return err
			},
		}),
	)

	// Act
	res := ckr.Check(context.Background())

	// Assert
	require.NotNil(t, res.Details)
	assert.Equal(t, expectedStatus, res.Status)

	for _, checkName := range []string{"check1", "check2"} {
		_, checkResultExists := res.Details[checkName]
		assert.True(t, checkResultExists)
	}
}

func TestWhenChecksExecutedThenAggregatedResultUp(t *testing.T) {
	t.Parallel()
	doTestCheckerCheckFunc(t, 0, nil, StatusUp)
}

func TestWhenOneCheckFailedThenAggregatedResultDown(t *testing.T) {
	t.Parallel()
	doTestCheckerCheckFunc(t, 0, fmt.Errorf("this is a check error"), StatusDown)
}

func TestCheckSuccessNotAllChecksExecutedYet(t *testing.T) {
	t.Parallel()
	doTestCheckerCheckFunc(t, 5*time.Hour, nil, StatusUnknown)
}

func TestPanicRecovery(t *testing.T) {
	t.Parallel()

	// Arrange
	expectedPanicMsg := "test message"
	ckr := NewChecker(
		WithCheck(Check{
			Name: "iPanic",
			Check: func(context.Context) error {
				panic(expectedPanicMsg)
			},
		}),
	)

	// Act
	res := ckr.Check(context.Background())

	// Assert
	require.NotNil(t, res.Details)
	assert.Equal(t, StatusDown, res.Status)

	checkRes, checkResultExists := res.Details["iPanic"]
	assert.True(t, checkResultExists)
	require.Error(t, checkRes.Error)
	assert.Equal(t, expectedPanicMsg, (checkRes.Error).Error())
}
