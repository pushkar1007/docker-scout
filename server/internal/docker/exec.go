package docker

import (
	"context"

	"github.com/moby/moby/client"
)

type ExecConfig struct {
	Cmd     []string
	TTY     bool
	Env     []string
	Workdir string
	User    string
	Cols    uint
	Rows    uint
}

func CreateExec(ctx context.Context, cli *client.Client, containerID string, cfg ExecConfig) (string, error) {
	res, err := cli.ExecCreate(ctx, containerID, client.ExecCreateOptions{
		TTY:          cfg.TTY,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Env:          cfg.Env,
		WorkingDir:   cfg.Workdir,
		User:         cfg.User,
		Cmd:          cfg.Cmd,
		ConsoleSize: client.ConsoleSize{
			Height: cfg.Rows,
			Width:  cfg.Cols,
		},
	})
	if err != nil {
		return "", err
	}
	return res.ID, nil
}

func AttachExec(ctx context.Context, cli *client.Client, execID string, tty bool, cols, rows uint) (client.HijackedResponse, error) {
	res, err := cli.ExecAttach(ctx, execID, client.ExecAttachOptions{
		TTY: tty,
		ConsoleSize: client.ConsoleSize{
			Height: rows,
			Width:  cols,
		},
	})
	if err != nil {
		return client.HijackedResponse{}, err
	}
	return res.HijackedResponse, nil
}

func ResizeExec(ctx context.Context, cli *client.Client, execID string, cols, rows uint) error {
	_, err := cli.ExecResize(ctx, execID, client.ExecResizeOptions{
		Height: rows,
		Width:  cols,
	})
	return err
}
