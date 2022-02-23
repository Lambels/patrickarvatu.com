package sqlite_test

import (
	"context"
	"reflect"
	"testing"

	pa "github.com/Lambels/patrickarvatu.com"
	"github.com/Lambels/patrickarvatu.com/sqlite"
)

// TODO: go with debugger in depth through each step to observe behaviour.

func TestCreateOrUpdateProject(t *testing.T) {
	t.Run("Ok Create or Update Call", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		projectService := sqlite.NewProjectService(db)

		backgroundCtx := context.Background()

		user := &pa.User{
			Name:    "Jhon Doe",
			Email:   "jhon@doe.com",
			IsAdmin: true,
		} // no need to create user as CreateOrUpdateProject doesent check any keys.

		adminUsrCtx := pa.NewContextWithUser(backgroundCtx, user)

		t.Run("step 1", func(t *testing.T) {
			// current projects: []
			pj1 := &pa.Project{
				Name:        "patrickarvatu.com",
				Description: "Some description",
				Topics: []string{
					"go",
				},
				HtmlURL: "https://github.com/Lambels/patrickarvatu.com",
			}

			// create project.
			if err := projectService.CreateOrUpdateProject(adminUsrCtx, pj1); err != nil {
				t.Fatal(err)
			} else if pj1.ID == 0 {
				t.Fatal("got id = 0")
			}

			// assert creation.
			if gotProject, err := projectService.FindProjectByID(backgroundCtx, pj1.ID); err != nil {
				t.Fatal(err)
			} else if !reflect.DeepEqual(pj1, gotProject) {
				t.Fatal("DeepEqual: gotProject != pj1")
			}
		})

		t.Run("step 2", func(t *testing.T) {
			// current projects: [patrickarvatu.com]

			pj1 := &pa.Project{
				Name:        "patrickarvatu.com",
				Description: "Some other description",
				Topics: []string{
					"go",
					"web dev",
					"js",
					"golang",
				},
				HtmlURL: "https://github.com/Lambels/patrickarvatu.com",
			}

			// update project 1.
			if err := projectService.CreateOrUpdateProject(adminUsrCtx, pj1); err != nil {
				t.Fatal(err)
			}

			// for reflect.DeepEqual test.
			pj1.ID = 1

			// assert update.
			if gotProject, err := projectService.FindProjectByID(backgroundCtx, 1); err != nil {
				t.Fatal(err)
			} else if !reflect.DeepEqual(pj1, gotProject) {
				t.Fatal("DeepEqual: gotProject != pj1")
			}
		})

		t.Run("step 3", func(t *testing.T) {
			// current projects: [patrickarvatu.com]
			pj1 := &pa.Project{
				Name:        "Lambels",
				Description: "No description :(",
				HtmlURL:     "https://github.com/Lambels/Lambels",
			}

			// create project.
			if err := projectService.CreateOrUpdateProject(adminUsrCtx, pj1); err != nil {
				t.Fatal(err)
			} else if pj1.ID == 0 {
				t.Fatal("got id = 0")
			}

			// assert creation.
			if gotProject, err := projectService.FindProjectByID(backgroundCtx, pj1.ID); err != nil {
				t.Fatal(err)
			} else if !reflect.DeepEqual(pj1, gotProject) {
				t.Fatal("DeepEqual: gotProject != pj1")
			}
		})

		t.Run("step 4", func(t *testing.T) {
			// current projects: [patrickarvatu.com, Lambels]

			pj1, err := projectService.FindProjectByID(backgroundCtx, 1)
			if err != nil {
				t.Fatal(err)
			}

			pj2, err := projectService.FindProjectByID(backgroundCtx, 2)
			if err != nil {
				t.Fatal(err)
			}

			// assert previous steps.
			if projects, n, err := projectService.FindProjects(backgroundCtx, pa.ProjectFilter{}); err != nil {
				t.Fatal(err)
			} else if n != 2 {
				t.Fatal("n != 2")
			} else if len(projects) != 2 {
				t.Fatal("len(projects) != 2")
			} else if !reflect.DeepEqual(projects[0], pj1) {
				t.Fatal("DeepEqual: projects[0] != pj1")
			} else if !reflect.DeepEqual(projects[1], pj2) {
				t.Fatal("DeepEqual: projects[1] != pj2")
			}
		})
	})

	t.Run("Bad Create Or Update Call (Un Auth)", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		projectService := sqlite.NewProjectService(db)

		backgroundCtx := context.Background()

		usrCtx := pa.NewContextWithUser(backgroundCtx, &pa.User{
			Name:  "jhon DOE",
			Email: "jhon@doe.com",
		}) // no need to create user as CreateProject doesent check any keys.

		adminUsrCtx := pa.NewContextWithUser(backgroundCtx, &pa.User{
			IsAdmin: true,
		})

		t.Run("Create", func(t *testing.T) {
			project := &pa.Project{
				Name:    "idk",
				HtmlURL: "idk",
			}

			// create project (Un Auth).
			if err := projectService.CreateOrUpdateProject(usrCtx, project); pa.ErrorCode(err) != pa.EUNAUTHORIZED {
				t.Fatal("expected UnAuth error")
			} else if project.ID != 0 {
				t.Fatal("got id != 0")
			}
		})

		t.Run("Update", func(t *testing.T) {
			project := &pa.Project{
				Name:    "idk",
				HtmlURL: "idk",
			}

			// create project.
			MustCreateOrUpdateProject(t, db, adminUsrCtx, project)

			// update project (Un Auth).
			if err := projectService.CreateOrUpdateProject(usrCtx, project); pa.ErrorCode(err) != pa.EUNAUTHORIZED {
				t.Fatal("expected UnAuth error")
			}
		})
	})
}

func TestDeleteProject(t *testing.T) {
	t.Run("Ok Delete Call", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		projectService := sqlite.NewProjectService(db)

		backgroundCtx := context.Background()

		adminUsrCtx := pa.NewContextWithUser(backgroundCtx, &pa.User{
			IsAdmin: true,
		})

		project := &pa.Project{
			Name:    "idk",
			HtmlURL: "idk",
		}

		// create project.
		MustCreateOrUpdateProject(t, db, adminUsrCtx, project)

		// delete project.
		if err := projectService.DeleteProject(adminUsrCtx, project.Name); err != nil {
			t.Fatal(err)
		}

		// assert deletion.
		if _, err := projectService.FindProjectByID(backgroundCtx, project.ID); pa.ErrorCode(err) != pa.ENOTFOUND {
			t.Fatal("err != ENOTFOUND")
		}
	})

	t.Run("Bad Delete Call (Un Auth)", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		projectService := sqlite.NewProjectService(db)

		user := &pa.User{
			Name:    "Jhon Doe",
			Email:   "jhon@doe.com",
			IsAdmin: true,
		} // no need to create user as DeleteProject doesent check any keys.

		user2 := &pa.User{
			Name:  "Lambels",
			Email: "Lamb@Lambels.com",
		}

		adminUsrCtx := pa.NewContextWithUser(backgroundCtx, user)
		usr2Ctx := pa.NewContextWithUser(backgroundCtx, user2)

		project := &pa.Project{
			Name:    "idk",
			HtmlURL: "idk",
		}

		// create project.
		MustCreateOrUpdateProject(t, db, adminUsrCtx, project)

		// delete project (Un Auth).
		if err := projectService.DeleteProject(usr2Ctx, "idk"); pa.ErrorCode(err) != pa.EUNAUTHORIZED {
			t.Fatal("err != EUNAUTHORIZED")
		}
	})

	t.Run("Bad Delete Call (Not Found)", func(t *testing.T) {
		db := MustOpenTempDB(t)
		defer MustCloseDB(t, db)

		backgroundCtx := context.Background()

		projectService := sqlite.NewProjectService(db)

		user := &pa.User{
			Name:    "Jhon Doe",
			Email:   "jhon@doe.com",
			IsAdmin: true,
		} // no need to create user as DeleteProject doesent check any keys.

		adminUsrCtx := pa.NewContextWithUser(backgroundCtx, user)

		// delete project (Not Found).
		if err := projectService.DeleteProject(adminUsrCtx, "fdsfsdv"); pa.ErrorCode(err) != pa.ENOTFOUND {
			t.Fatal("err != ENOTFOUND")
		}
	})
}

func MustCreateOrUpdateProject(t *testing.T, db *sqlite.DB, ctx context.Context, project *pa.Project) {
	t.Helper()
	if err := sqlite.NewProjectService(db).CreateOrUpdateProject(ctx, project); err != nil {
		t.Fatal(err)
	}
}
