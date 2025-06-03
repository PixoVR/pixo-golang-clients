package platform_test

import (
	"context"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/go-faker/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("SuccessFactors", func() {
	var (
		ctx = context.Background()
	)

	Context("Learning History", func() {

		var (
			session *platform.Session
		)

		BeforeEach(func() {
			uuid := faker.UUIDHyphenated()
			session = &platform.Session{
				ModuleID: moduleID,
				OrgID:    orgID,
				UUID:     &uuid,
			}

			Expect(tokenClient.CreateSession(ctx, session)).To(Succeed())

			Expect(session).NotTo(BeNil())
			Expect(session.ID).NotTo(BeZero())

			session.CompletedAt = time.Now().Add(time.Minute * 5).Format(time.RFC3339)
			session.Status = "completed"
			session.RawScore = 2342

			_, err := tokenClient.UpdateSession(ctx, *session)
			Expect(err).NotTo(HaveOccurred())
		})

		It("can get learning history records by date range and org ID", func() {
			startDate := time.Now().Add(-time.Minute)
			EndDate := time.Now().Add(time.Minute * 10)
			params := platform.LearningHistoryParams{
				OrgID:     orgID,
				StartDate: startDate,
				EndDate:   EndDate,
			}

			learningHistoryRecords, err := tokenClient.GetLearningHistoryRecords(ctx, params)
			Expect(err).NotTo(HaveOccurred())
			Expect(learningHistoryRecords).NotTo(BeNil())
			Expect(len(learningHistoryRecords)).To(BeNumerically(">", 0))

			var foundRecord platform.LearningHistory
			for _, record := range learningHistoryRecords {
				if record.ID != session.ID {
					continue
				}

				foundRecord = record
				break
			}

			Expect(foundRecord).NotTo(BeNil(), "learning history record not found")
			Expect(foundRecord.UserID).To(Equal(session.UserID))
			Expect(foundRecord.ModuleID).To(Equal(session.ModuleID))
			Expect(foundRecord.Score).To(Equal(session.RawScore))
			Expect(foundRecord.SessionDuration).To(BeNumerically("~", 300, 10))
		})
	})

	Context("Course Data", func() {
		It("can get course data by org ID", func() {
			courseDataRecords, err := tokenClient.GetCourseDataRecords(ctx, orgID)
			Expect(err).NotTo(HaveOccurred())

			Expect(courseDataRecords).NotTo(BeNil())
			Expect(len(courseDataRecords)).To(BeNumerically(">", 0))

			var foundRecord platform.CourseData
			for _, record := range courseDataRecords {
				if record.ID != moduleID {
					continue
				}

				foundRecord = record
				break
			}

			Expect(foundRecord).NotTo(BeNil(), "course data record not found for seeded module")
			Expect(foundRecord.Name).To(ContainSubstring("360"))
			Expect(foundRecord.Description).To(Equal("PIXO VR 360 Video Player "))
		})
	})

	Context("Org Success Factors", func() {
		It("can get org success factors", func() {
			orgSuccessFactors, err := tokenClient.GetOrgSuccessFactors(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(orgSuccessFactors).NotTo(BeNil())
		})
	})

})
