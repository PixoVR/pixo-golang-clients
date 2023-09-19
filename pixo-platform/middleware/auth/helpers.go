package auth

import (
	"errors"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func GetCurrentUser(c *gin.Context) (*platform.User, error) {
	context := GetContext(c.Request.Context())

	tokenString := ExtractSecretKey(c.Request)

	if tokenString == "" {
		log.Debug().Msg("token not found")
		return nil, errors.New("token not found")
	}

	rawToken, err := ParseJWT(tokenString)
	if err != nil {
		log.Debug().Err(err).Msg("error parsing JWT")
		return nil, err
	}

	if context.FindUserByID == nil {
		log.Debug().Msg("user not found")
		return nil, errors.New("user not found")
	}

	user, err := context.FindUserByID(rawToken.UserID)
	if err != nil {
		log.Debug().Msg("user not found")
		return nil, err
	}

	return user, nil
}

//func GetOrgAccess(orgService services.OrgService, userOrg *models.Org) []*models.Org {
//	if userOrg == nil {
//		return nil
//	}
//
//	if userOrg.Type == "platform" {
//		allOrgs, err := orgService.GetAll()
//		if err != nil {
//			log.Debug().Msg(ErrOrgAccess)
//			return nil
//		}
//		return allOrgs
//	}
//
//	orgsInTree, err := orgService.GetOrgsInTree(userOrg)
//	if err != nil {
//		log.Debug().Msg(ErrOrgAccess)
//		return nil
//	}
//	return orgsInTree
//}
