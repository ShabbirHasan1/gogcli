package cmd

import "github.com/steipete/gogcli/internal/googleapi"

var newClassroomService = googleapi.NewClassroom

type ClassroomCmd struct {
	Courses         ClassroomCoursesCmd         `cmd:"" help:"Courses"`
	Students        ClassroomStudentsCmd        `cmd:"" help:"Course students"`
	Teachers        ClassroomTeachersCmd        `cmd:"" help:"Course teachers"`
	Coursework      ClassroomCourseworkCmd      `cmd:"" name:"coursework" help:"Coursework"`
	Submissions     ClassroomSubmissionsCmd     `cmd:"" help:"Student submissions"`
	Announcements   ClassroomAnnouncementsCmd   `cmd:"" help:"Announcements"`
	Topics          ClassroomTopicsCmd          `cmd:"" help:"Topics"`
	Invitations     ClassroomInvitationsCmd     `cmd:"" help:"Invitations"`
	Guardians       ClassroomGuardiansCmd       `cmd:"" help:"Guardians"`
	GuardianInvites ClassroomGuardianInvitesCmd `cmd:"" name:"guardian-invitations" help:"Guardian invitations"`
	Profile         ClassroomProfileCmd         `cmd:"" help:"User profiles"`
}
