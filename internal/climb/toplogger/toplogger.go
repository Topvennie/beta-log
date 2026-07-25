package toplogger

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Topvennie/beta-log/internal/database/model"
	"github.com/Topvennie/beta-log/internal/database/repository"
)

const (
	baseURL      = "https://app.toplogger.nu"
	uploadURL    = "https://upload.toplogger.nu"
	queryDayList = `[{"operationName":"climbDaysSessionsList","variables":{"pagination":{"page":%d,"perPage":10},"userId":"%s"},"query":"query climbDaysSessionsList($userId: ID!, $bouldersTotalTriesMin: Int, $routesTotalTriesMin: Int, $statsAtDateMin: DateTime, $statsAtDateMax: DateTime, $pagination: PaginationInputClimbDays) {\n  climbDaysPaginated(\n    userId: $userId\n    totalTriesMin: 1\n    bouldersTotalTriesMin: $bouldersTotalTriesMin\n    routesTotalTriesMin: $routesTotalTriesMin\n    statsAtDateMin: $statsAtDateMin\n    statsAtDateMax: $statsAtDateMax\n    pagination: $pagination\n    updateDayStatsIfOld: true\n  ) {\n    pagination {\n      ...pagination\n      __typename\n    }\n    data {\n      id\n      ...climbDayForSessionsList\n      __typename\n    }\n    __typename\n  }\n}\n\nfragment climbDayForUseSessionSummaryTitle on ClimbDay {\n  id\n  title\n  routesTotalTries\n  bouldersTotalTries\n  routesDayGradeMax\n  bouldersDayGradeMax\n  __typename\n}\n\nfragment climbDayForSessionSummaryTitle on ClimbDay {\n  id\n  description\n  statsAtDate\n  bouldersDayGrade\n  routesDayGrade\n  gym {\n    id\n    name\n    nameSlug\n    iconPath\n    __typename\n  }\n  user {\n    id\n    fullName\n    avatarUploadPath\n    __typename\n  }\n  ...climbDayForUseSessionSummaryTitle\n  __typename\n}\n\nfragment gymForGradingSystem on Gym {\n  id\n  gradingSystemRoutes\n  gradingSystemBoulders\n  gradingSystemRoutesCustom\n  gradingSystemBouldersCustom\n  __typename\n}\n\nfragment climbDayForSessionSummaryMetrics on ClimbDay {\n  id\n  bouldersTotalTries\n  bouldersDayGrade\n  bouldersDayGradeFlPct\n  bouldersDayGradeRepeatPct\n  routesTotalTries\n  routesDayGrade\n  routesDayGradeOsPct\n  routesDayGradeRepeatPct\n  routesTotalHeight\n  gym {\n    ...gymForGradingSystem\n    __typename\n  }\n  __typename\n}\n\nfragment gymForClimbTagColor on Gym {\n  id\n  climbGroups {\n    id\n    climbGroupBy\n    color\n    __typename\n  }\n  __typename\n}\n\nfragment gymForSimpleMapClimb on Gym {\n  id\n  ...gymForGradingSystem\n  ...gymForClimbTagColor\n  __typename\n}\n\nfragment climbForClimbTagColor on Climb {\n  id\n  climbGroupClimbs {\n    id\n    climbGroupId\n    __typename\n  }\n  __typename\n}\n\nfragment climbForSimpleMapClimb on Climb {\n  id\n  positionX\n  positionY\n  grade\n  label\n  climbType\n  holdColor {\n    id\n    color\n    colorSecondary\n    __typename\n  }\n  ...climbForClimbTagColor\n  __typename\n}\n\nfragment climbDayForSessionMap on ClimbDay {\n  id\n  gym {\n    id\n    floorplanPath\n    ...gymForSimpleMapClimb\n    __typename\n  }\n  climbUserDaysRoutes: climbUserDays(climbType: \"route\", limit: 5) {\n    id\n    tickType\n    wasRepeat\n    climb {\n      id\n      ...climbForSimpleMapClimb\n      __typename\n    }\n    __typename\n  }\n  climbUserDaysBoulders: climbUserDays(climbType: \"boulder\", limit: 10) {\n    id\n    tickType\n    wasRepeat\n    climb {\n      id\n      ...climbForSimpleMapClimb\n      __typename\n    }\n    __typename\n  }\n  __typename\n}\n\nfragment climbDayForSessionSummary on ClimbDay {\n  id\n  statsAtDate\n  bouldersTotalTries\n  bouldersDayGradeMax\n  routesTotalTries\n  routesDayGradeMax\n  ...climbDayForSessionSummaryTitle\n  ...climbDayForSessionSummaryMetrics\n  ...climbDayForSessionMap\n  __typename\n}\n\nfragment climbDayForSessionRoute on ClimbDay {\n  id\n  userId\n  __typename\n}\n\nfragment climbDayForFeedback on ClimbDay {\n  id\n  userId\n  likesCount\n  commentsCount\n  likeMe {\n    id\n    __typename\n  }\n  likesForAvatarStack: comments(type: \"LIKE\", limit: 3) {\n    id\n    user {\n      id\n      avatarUploadPath\n      __typename\n    }\n    __typename\n  }\n  ...climbDayForSessionRoute\n  __typename\n}\n\nfragment climbDayForLikeBtn on ClimbDay {\n  id\n  userId\n  likeMe {\n    id\n    __typename\n  }\n  ...climbDayForFeedback\n  ...climbDayForSessionRoute\n  __typename\n}\n\nfragment climbDayForSession on ClimbDay {\n  id\n  ...climbDayForSessionSummary\n  ...climbDayForSessionRoute\n  ...climbDayForFeedback\n  ...climbDayForLikeBtn\n  __typename\n}\n\nfragment pagination on Pagination {\n  total\n  page\n  perPage\n  orderBy {\n    key\n    order\n    __typename\n  }\n  __typename\n}\n\nfragment climbDayForSessionsList on ClimbDay {\n  id\n  ...climbDayForSession\n  __typename\n}"}]`
	queryDayLog  = `[{"operationName":"climbLogsSession","variables":{"pagination":{"orderBy":[{"key":"points","order":"desc"}]},"gymId":"%s","userId":"%s","climbedAtDate":"%s","climbType":"boulder"},"query":"query climbLogsSession($gymId: ID, $userId: ID!, $climbedAtDate: DateTime, $pagination: PaginationInputClimbLogs, $compRoundId: ID, $climbType: ClimbType) {\n  climbLogs(\n    gymId: $gymId\n    userId: $userId\n    climbedAtDate: $climbedAtDate\n    climbType: $climbType\n    pagination: $pagination\n  ) {\n    pagination {\n      ...pagination\n      __typename\n    }\n    data {\n      id\n      climbId\n      ...climbLogForSessionClimb\n      __typename\n    }\n    __typename\n  }\n}\n\nfragment climbForClimbLog on Climb {\n  id\n  leadRequired\n  outAt\n  compRoundClimb(compRoundId: $compRoundId) {\n    id\n    leadRequired\n    __typename\n  }\n  __typename\n}\n\nfragment climbLogForRemove on ClimbLog {\n  id\n  gymId\n  climbId\n  __typename\n}\n\nfragment compClimbLogForScoreSystemResultsPoints on CompClimbLog {\n  id\n  points\n  pointsBase\n  pointsBonus\n  pointsJson\n  __typename\n}\n\nfragment compClimbLogForScoreSystemResults on CompClimbLog {\n  id\n  ...compClimbLogForScoreSystemResultsPoints\n  __typename\n}\n\nfragment climbLog on ClimbLog {\n  id\n  gymId\n  userId\n  climbId\n  climbType\n  topped\n  foreknowledge\n  zones\n  clips\n  holds\n  duration\n  lead\n  hangs\n  comments\n  tryIndex\n  tickIndex\n  ticked\n  tickType\n  points\n  climbedAtDate\n  ...climbLogForRemove\n  compClimbLog(compRoundId: $compRoundId) {\n    id\n    points\n    pointsBase\n    pointsBonus\n    pointsJson\n    ...compClimbLogForScoreSystemResults\n    __typename\n  }\n  __typename\n}\n\nfragment pagination on Pagination {\n  total\n  page\n  perPage\n  orderBy {\n    key\n    order\n    __typename\n  }\n  __typename\n}\n\nfragment climbLogForSessionClimb on ClimbLog {\n  id\n  userId\n  points\n  pointsBonus\n  tryIndex\n  tickIndex\n  ticked\n  tickType\n  climb {\n    id\n    name\n    grade\n    climbType\n    leadRequired\n    outAt\n    holdColor {\n      id\n      color\n      colorSecondary\n      __typename\n    }\n    wall {\n      id\n      nameLoc\n      __typename\n    }\n    gym {\n      id\n      name\n      nameSlug\n      __typename\n    }\n    ...climbForClimbLog\n    __typename\n  }\n  ...climbLog\n  __typename\n}"}]`
	queryRefetch = `[{"operationName":"authSigninRefreshToken","variables":{"refreshToken":"%s"},"query":"mutation authSigninRefreshToken($refreshToken: JWT!) {\n  tokens: authSigninRefreshToken(refreshToken: $refreshToken) {\n    ...authTokens\n    __typename\n  }\n}\n\nfragment authTokens on AuthTokens {\n  access {\n    token\n    expiresAt\n    __typename\n  }\n  refresh {\n    token\n    expiresAt\n    __typename\n  }\n  __typename\n}"}]`
)

var (
	ErrUnauthorized    = errors.New("unauthorized")
	ErrExpired         = errors.New("tokens expired")
	ErrNoTokenResponse = errors.New("no new tokens received")
)

type Client struct {
	day     repository.ClimbDay
	setting repository.Setting
}

func New() *Client {
	return &Client{
		day:     *repository.NewClimbDay(),
		setting: *repository.NewSetting(),
	}
}

func (c *Client) Fetch(ctx context.Context, user model.User) ([]model.ClimbDay, error) {
	setting, err := c.setting.GetByUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if setting.ClimbToploggerUserID == "" || setting.ClimbToploggerAuthToken == "" || setting.ClimbToploggerRefreshToken == "" || setting.ClimbTopLoggerExpiration.IsZero() {
		return nil, nil
	}
	if setting.ClimbTopLoggerExpiration.Before(time.Now()) {
		// Refresh token expired
		// Remove the tokens from the settings
		if err := c.resetSetting(ctx, *setting); err != nil {
			return nil, err
		}
		return nil, ErrExpired
	}

	// Refresh tokens
	tokens, err := c.Refresh(ctx, *setting)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) || errors.Is(err, ErrNoTokenResponse) {
			if errReset := c.resetSetting(ctx, *setting); errReset != nil {
				return nil, errReset
			}
			return nil, err
		}
		return nil, err
	}

	setting.ClimbToploggerAuthToken = tokens.Access.Token
	setting.ClimbToploggerRefreshToken = tokens.Refresh.Token
	setting.ClimbTopLoggerExpiration = tokens.Refresh.ExpiresAt

	if err := c.setting.ToploggerUpdate(ctx, *setting); err != nil {
		return nil, err
	}

	// Get all days
	climbDays, err := c.getDays(ctx, *setting)
	if err != nil {
		return nil, err
	}

	// Process each day
	gymMap := map[string]model.ClimbGym{}
	days := make([]model.ClimbDay, 0, len(climbDays))

	for _, day := range climbDays {
		date, err := parseDate(day.StatsAtDate)
		if err != nil {
			return nil, err
		}

		gym, ok := gymMap[day.Gym.ID]
		if !ok {
			gym = model.ClimbGym{
				UserID:     user.ID,
				ExternalID: day.Gym.ID,
				Name:       day.Gym.Name,
				IconPath:   fmt.Sprintf("%s/%s", uploadURL, day.Gym.IconPath),
				Source:     model.ClimbSourceToplogger,
			}
			gymMap[day.Gym.ID] = gym
		}

		dayClimbs, err := c.getDayClimbs(ctx, *setting, day)
		if err != nil {
			return nil, err
		}

		climbs := make([]model.Climb, 0, len(dayClimbs))
		for _, climb := range dayClimbs {
			if climb.TickIndex < 0 {
				continue
			}

			finishType := model.ClimbFinishFlash
			if climb.TickType == 1 {
				finishType = model.ClimbFinishTop
			}
			if climb.TickIndex > 0 {
				finishType = model.ClimbFinishRepeat
			}

			climbs = append(climbs, model.Climb{
				UserID:     user.ID,
				ExternalID: climb.Climb.ID,
				Grade:      climb.Climb.Grade,
				HoldColor:  climb.Climb.HoldColor.Color,
				ClimbType:  model.ClimbType(climb.Climb.ClimbType),
				FinishType: finishType,
				Source:     model.ClimbSourceToplogger,
			})
		}

		days = append(days, model.ClimbDay{
			UserID:     user.ID,
			ExternalID: day.ID,
			Date:       date,
			Gym:        gym,
			Climbs:     climbs,
			Source:     model.ClimbSourceToplogger,
		})

		// Small timeout
		time.Sleep(500 * time.Millisecond)
	}

	return days, nil
}

// getDays get all climbing days taking pagination into account
func (c *Client) getDays(ctx context.Context, setting model.Setting) ([]climbDay, error) {
	var climbDays []climbDay

	page := 1
	for {
		climbDayPaginated, err := c.fetchDayListPage(ctx, setting, page)
		if err != nil {
			if errors.Is(err, ErrUnauthorized) {
				if err := c.resetSetting(ctx, setting); err != nil {
					return nil, err
				}
				return nil, ErrUnauthorized
			}

			return nil, err
		}

		climbDays = append(climbDays, climbDayPaginated.Data...)

		if climbDayPaginated.Pagination.Page*climbDayPaginated.Pagination.PerPage >= climbDayPaginated.Pagination.Total {
			break
		}

		page++
	}

	return climbDays, nil
}

func (c *Client) getDayClimbs(ctx context.Context, setting model.Setting, day climbDay) ([]climbLog, error) {
	var climbLogs []climbLog

	page := 1
	for {
		climbLogPaginated, err := c.fetchDayLogPage(ctx, setting, day, page)
		if err != nil {
			if errors.Is(err, ErrUnauthorized) {
				if err := c.resetSetting(ctx, setting); err != nil {
					return nil, err
				}
				return nil, ErrUnauthorized
			}

			return nil, err
		}

		climbLogs = append(climbLogs, climbLogPaginated.Data...)

		if climbLogPaginated.Pagination.Page*climbLogPaginated.Pagination.PerPage >= climbLogPaginated.Pagination.Total {
			break
		}

		page++
	}

	return climbLogs, nil
}

func (c *Client) fetchDayListPage(ctx context.Context, setting model.Setting, page int) (climbDayPaginated, error) {
	query := fmt.Sprintf(queryDayList, page, setting.ClimbToploggerUserID)

	type climbResponse struct {
		Data struct {
			ClimbDayPaginated climbDayPaginated `json:"climbDaysPaginated"`
		} `json:"data"`
	}
	var climbResult []climbResponse

	respBody, code, err := c.request(ctx, setting.ClimbToploggerAuthToken, http.MethodPost, "graphql", strings.NewReader(query))
	if err != nil {
		return climbDayPaginated{}, err
	}
	if code != 200 {
		return climbDayPaginated{}, fmt.Errorf("wrong status code %d", code)
	}

	if err := json.Unmarshal(respBody, &climbResult); err != nil {
		return climbDayPaginated{}, fmt.Errorf("unmarshal response body %w", err)
	}

	if len(climbResult) == 0 || len(climbResult[0].Data.ClimbDayPaginated.Data) == 0 {
		// Maybe it was an error
		if err := getError(respBody); err != nil {
			return climbDayPaginated{}, err
		}

		// Nevermind no error, it just doesn't contain any data
		return climbDayPaginated{}, nil
	}

	return climbResult[0].Data.ClimbDayPaginated, nil
}

func (c *Client) fetchDayLogPage(ctx context.Context, setting model.Setting, day climbDay, page int) (climbLogPaginated, error) {
	query := fmt.Sprintf(queryDayLog, day.Gym.ID, setting.ClimbToploggerUserID, day.StatsAtDate)

	type climbResponse struct {
		Data struct {
			ClimbLogs climbLogPaginated `json:"climbLogs"`
		} `json:"data"`
	}
	var climbResult []climbResponse

	respBody, code, err := c.request(ctx, setting.ClimbToploggerAuthToken, http.MethodPost, "graphql", strings.NewReader(query))
	if err != nil {
		return climbLogPaginated{}, err
	}
	if code != 200 {
		return climbLogPaginated{}, fmt.Errorf("wrong status code %d", code)
	}

	if err := json.Unmarshal(respBody, &climbResult); err != nil {
		return climbLogPaginated{}, fmt.Errorf("unmarshal response body %w", err)
	}

	if len(climbResult) == 0 || len(climbResult[0].Data.ClimbLogs.Data) == 0 {
		// Maybe it was an error
		if err := getError(respBody); err != nil {
			return climbLogPaginated{}, err
		}

		// Nevermind no error, it just doesn't contain any data
		return climbLogPaginated{}, nil
	}

	return climbResult[0].Data.ClimbLogs, nil
}

func (c *Client) Refresh(ctx context.Context, setting model.Setting) (Token, error) {
	query := fmt.Sprintf(queryRefetch, setting.ClimbToploggerRefreshToken)

	type response struct {
		Data struct {
			Tokens Token `json:"tokens"`
		} `json:"data"`
	}
	var result []response

	respBody, code, err := c.request(ctx, setting.ClimbToploggerRefreshToken, http.MethodPost, "graphql", strings.NewReader(query))
	if err != nil {
		return Token{}, err
	}
	if code != 200 {
		if code == 400 {
			return Token{}, ErrUnauthorized
		}
		return Token{}, fmt.Errorf("wrong status code %d", code)
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return Token{}, fmt.Errorf("unmarshal response body %w", err)
	}

	if len(result) == 0 || result[0].Data.Tokens.Access.Token == "" || result[0].Data.Tokens.Refresh.Token == "" || result[0].Data.Tokens.Refresh.ExpiresAt.IsZero() {
		// No response
		// Maybe it is an error
		if err := getError(respBody); err != nil {
			return Token{}, err
		}

		// Nevermind no error
		return Token{}, ErrNoTokenResponse
	}

	return result[0].Data.Tokens, nil
}
