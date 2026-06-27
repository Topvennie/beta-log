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
	queryDay     = `[{"operationName":"climbDaysSessionsList","variables":{"pagination":{"page":%d,"perPage":10},"userId":"%s"},"query":"query climbDaysSessionsList($userId: ID!, $bouldersTotalTriesMin: Int, $routesTotalTriesMin: Int, $statsAtDateMin: DateTime, $statsAtDateMax: DateTime, $pagination: PaginationInputClimbDays) {\n  climbDaysPaginated(\n    userId: $userId\n    totalTriesMin: 1\n    bouldersTotalTriesMin: $bouldersTotalTriesMin\n    routesTotalTriesMin: $routesTotalTriesMin\n    statsAtDateMin: $statsAtDateMin\n    statsAtDateMax: $statsAtDateMax\n    pagination: $pagination\n    updateDayStatsIfOld: true\n  ) {\n    pagination {\n      ...pagination\n      __typename\n    }\n    data {\n      id\n      ...climbDayForSessionsList\n      __typename\n    }\n    __typename\n  }\n}\n\nfragment climbDayForUseSessionSummaryTitle on ClimbDay {\n  id\n  title\n  routesTotalTries\n  bouldersTotalTries\n  routesDayGradeMax\n  bouldersDayGradeMax\n  __typename\n}\n\nfragment climbDayForSessionSummaryTitle on ClimbDay {\n  id\n  description\n  statsAtDate\n  bouldersDayGrade\n  routesDayGrade\n  gym {\n    id\n    name\n    nameSlug\n    iconPath\n    __typename\n  }\n  user {\n    id\n    fullName\n    avatarUploadPath\n    __typename\n  }\n  ...climbDayForUseSessionSummaryTitle\n  __typename\n}\n\nfragment gymForGradingSystem on Gym {\n  id\n  gradingSystemRoutes\n  gradingSystemBoulders\n  gradingSystemRoutesCustom\n  gradingSystemBouldersCustom\n  __typename\n}\n\nfragment climbDayForSessionSummaryMetrics on ClimbDay {\n  id\n  bouldersTotalTries\n  bouldersDayGrade\n  bouldersDayGradeFlPct\n  bouldersDayGradeRepeatPct\n  routesTotalTries\n  routesDayGrade\n  routesDayGradeOsPct\n  routesDayGradeRepeatPct\n  routesTotalHeight\n  gym {\n    ...gymForGradingSystem\n    __typename\n  }\n  __typename\n}\n\nfragment gymForClimbTagColor on Gym {\n  id\n  climbGroups {\n    id\n    climbGroupBy\n    color\n    __typename\n  }\n  __typename\n}\n\nfragment gymForSimpleMapClimb on Gym {\n  id\n  ...gymForGradingSystem\n  ...gymForClimbTagColor\n  __typename\n}\n\nfragment climbForClimbTagColor on Climb {\n  id\n  climbGroupClimbs {\n    id\n    climbGroupId\n    __typename\n  }\n  __typename\n}\n\nfragment climbForSimpleMapClimb on Climb {\n  id\n  positionX\n  positionY\n  grade\n  label\n  climbType\n  holdColor {\n    id\n    color\n    colorSecondary\n    __typename\n  }\n  ...climbForClimbTagColor\n  __typename\n}\n\nfragment climbDayForSessionMap on ClimbDay {\n  id\n  gym {\n    id\n    floorplanPath\n    ...gymForSimpleMapClimb\n    __typename\n  }\n  climbUserDaysRoutes: climbUserDays(climbType: \"route\", limit: 5) {\n    id\n    tickType\n    wasRepeat\n    climb {\n      id\n      ...climbForSimpleMapClimb\n      __typename\n    }\n    __typename\n  }\n  climbUserDaysBoulders: climbUserDays(climbType: \"boulder\", limit: 10) {\n    id\n    tickType\n    wasRepeat\n    climb {\n      id\n      ...climbForSimpleMapClimb\n      __typename\n    }\n    __typename\n  }\n  __typename\n}\n\nfragment climbDayForSessionSummary on ClimbDay {\n  id\n  statsAtDate\n  bouldersTotalTries\n  bouldersDayGradeMax\n  routesTotalTries\n  routesDayGradeMax\n  ...climbDayForSessionSummaryTitle\n  ...climbDayForSessionSummaryMetrics\n  ...climbDayForSessionMap\n  __typename\n}\n\nfragment climbDayForSessionRoute on ClimbDay {\n  id\n  userId\n  __typename\n}\n\nfragment climbDayForFeedback on ClimbDay {\n  id\n  userId\n  likesCount\n  commentsCount\n  likeMe {\n    id\n    __typename\n  }\n  likesForAvatarStack: comments(type: \"LIKE\", limit: 3) {\n    id\n    user {\n      id\n      avatarUploadPath\n      __typename\n    }\n    __typename\n  }\n  ...climbDayForSessionRoute\n  __typename\n}\n\nfragment climbDayForLikeBtn on ClimbDay {\n  id\n  userId\n  likeMe {\n    id\n    __typename\n  }\n  ...climbDayForFeedback\n  ...climbDayForSessionRoute\n  __typename\n}\n\nfragment climbDayForSession on ClimbDay {\n  id\n  ...climbDayForSessionSummary\n  ...climbDayForSessionRoute\n  ...climbDayForFeedback\n  ...climbDayForLikeBtn\n  __typename\n}\n\nfragment pagination on Pagination {\n  total\n  page\n  perPage\n  orderBy {\n    key\n    order\n    __typename\n  }\n  __typename\n}\n\nfragment climbDayForSessionsList on ClimbDay {\n  id\n  ...climbDayForSession\n  __typename\n}"}]`
	queryRefetch = `[{"operationName":"authSigninRefreshToken","variables":{"refreshToken":"%s"},"query":"mutation authSigninRefreshToken($refreshToken: JWT!) {\n  tokens: authSigninRefreshToken(refreshToken: $refreshToken) {\n    ...authTokens\n    __typename\n  }\n}\n\nfragment authTokens on AuthTokens {\n  access {\n    token\n    expiresAt\n    __typename\n  }\n  refresh {\n    token\n    expiresAt\n    __typename\n  }\n  __typename\n}"}]`
)

var errUnauthorized = errors.New("unauthorized")

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
		return nil, c.resetSetting(ctx, *setting)
	}

	// Refresh tokens
	tokens, err := c.refresh(ctx, *setting)
	if err != nil {
		if errors.Is(err, errUnauthorized) {
			return nil, c.resetSetting(ctx, *setting)
		}

		return nil, err
	}
	if tokens.Access.Token == "" || tokens.Refresh.Token == "" || tokens.Refresh.ExpiresAt.IsZero() {
		return nil, c.resetSetting(ctx, *setting)
	}

	setting.ClimbToploggerAuthToken = tokens.Access.Token
	setting.ClimbToploggerRefreshToken = tokens.Refresh.Token
	setting.ClimbTopLoggerExpiration = tokens.Refresh.ExpiresAt

	if err := c.setting.Update(ctx, *setting); err != nil {
		return nil, err
	}

	// Get all climb days taking the pagination into account
	var climbDays []climbDay

	page := 1
	for {
		climbDayPaginated, err := c.fetchPage(ctx, *setting, page)
		if err != nil {
			if errors.Is(err, errUnauthorized) {
				if err := c.resetSetting(ctx, *setting); err != nil {
					return nil, err
				}
			}

			return nil, err
		}

		climbDays = append(climbDays, climbDayPaginated.Data...)

		if climbDayPaginated.Pagination.Page*climbDayPaginated.Pagination.PerPage >= climbDayPaginated.Pagination.Total {
			break
		}

		page++
	}

	// Convert all days
	gymMap := map[string]model.ClimbGym{}
	gymColorsMap := map[string]string{}
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
			}
			gymMap[day.Gym.ID] = gym

			for _, climbGroup := range day.Gym.ClimbGroups {
				gymColorsMap[climbGroup.ID] = climbGroup.Color
			}
		}

		climbs := make([]model.Climb, 0, len(day.ClimbUserDaysBoulders))
		for _, climb := range day.ClimbUserDaysBoulders {
			finishType := model.ClimbFinishFlash
			if climb.TickType == 1 {
				finishType = model.ClimbFinishTop
			}
			if climb.WasRepeat {
				finishType = model.ClimbFinishRepeat
			}

			// Get first color
			color := ""
			if len(climb.Climb.ClimbGroupClimbs) > 0 {
				for _, climbGroupClimb := range climb.Climb.ClimbGroupClimbs {
					if c, ok := gymColorsMap[climbGroupClimb.ClimbGroupID]; ok {
						color = c
						break
					}
				}
			}

			climbs = append(climbs, model.Climb{
				UserID:     user.ID,
				ExternalID: climb.Climb.ID,
				Grade:      climb.Climb.Grade,
				Color:      color,
				HoldColor:  climb.Climb.HoldColor.Color,
				ClimbType:  model.ClimbType(climb.Climb.ClimbType),
				FinishType: finishType,
			})
		}

		days = append(days, model.ClimbDay{
			UserID:     user.ID,
			ExternalID: day.ID,
			Date:       date,
			Gym:        gym,
			Climbs:     climbs,
		})
	}

	return days, nil
}

func (c *Client) fetchPage(ctx context.Context, setting model.Setting, page int) (climbDayPaginated, error) {
	query := fmt.Sprintf(queryDay, page, setting.ClimbToploggerUserID)

	type climbResponse struct {
		Data struct {
			ClimbDayPaginated climbDayPaginated `json:"climbDaysPaginated"`
		} `json:"data"`
	}
	var climbResult []climbResponse

	respBody, err := c.request(ctx, setting.ClimbToploggerAuthToken, http.MethodPost, "graphql", strings.NewReader(query))
	if err != nil {
		return climbDayPaginated{}, err
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

func (c *Client) refresh(ctx context.Context, setting model.Setting) (token, error) {
	query := fmt.Sprintf(queryRefetch, setting.ClimbToploggerRefreshToken)

	type response struct {
		Data struct {
			Tokens token `json:"tokens"`
		} `json:"data"`
	}
	var result []response

	respBody, err := c.request(ctx, setting.ClimbToploggerRefreshToken, http.MethodPost, "graphql", strings.NewReader(query))
	if err != nil {
		return token{}, err
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return token{}, fmt.Errorf("unmarshal response body %w", err)
	}

	if len(result) == 0 || result[0].Data.Tokens.Access.Token == "" {
		// No response
		// Maybe it is an error
		if err := getError(respBody); err != nil {
			return token{}, err
		}

		// Nevermind no error
		return token{}, nil
	}

	return result[0].Data.Tokens, nil
}
