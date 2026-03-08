package cmd

import (
	"context"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/sheets/v4"
)

func requireCalendarService(ctx context.Context, flags *RootFlags) (string, *calendar.Service, error) {
	account, err := requireAccount(flags)
	if err != nil {
		return "", nil, err
	}
	svc, err := newCalendarService(ctx, account)
	if err != nil {
		return "", nil, err
	}
	return account, svc, nil
}

func requireDocsService(ctx context.Context, flags *RootFlags) (string, *docs.Service, error) {
	account, err := requireAccount(flags)
	if err != nil {
		return "", nil, err
	}
	svc, err := newDocsService(ctx, account)
	if err != nil {
		return "", nil, err
	}
	return account, svc, nil
}

func requireDriveService(ctx context.Context, flags *RootFlags) (string, *drive.Service, error) {
	account, err := requireAccount(flags)
	if err != nil {
		return "", nil, err
	}
	svc, err := newDriveService(ctx, account)
	if err != nil {
		return "", nil, err
	}
	return account, svc, nil
}

func requireGmailService(ctx context.Context, flags *RootFlags) (string, *gmail.Service, error) {
	account, err := requireAccount(flags)
	if err != nil {
		return "", nil, err
	}
	svc, err := newGmailService(ctx, account)
	if err != nil {
		return "", nil, err
	}
	return account, svc, nil
}

func requireSheetsService(ctx context.Context, flags *RootFlags) (string, *sheets.Service, error) {
	account, err := requireAccount(flags)
	if err != nil {
		return "", nil, err
	}
	svc, err := newSheetsService(ctx, account)
	if err != nil {
		return "", nil, err
	}
	return account, svc, nil
}
