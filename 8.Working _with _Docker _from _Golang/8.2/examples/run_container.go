package examples

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

func RunContainer() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation()) //собираем клиент с опциями из env переменных
	if err != nil {
		panic(err)
	}

	reader, err := cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{}) //стянуть изображение можно не только по имени но и по полному пути
	if err != nil {
		panic(err)
	}

	defer reader.Close()
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{ //запускаем контейнер с командой echo hello world
		Image: "alpine",
		Cmd:   []string{"echo", "hello world"},
		Tty:   false, //не содиняем вывод с консолью пользователя
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	} //запускаем контейнер

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning) //ждем пока контейнер не остановится
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true}) //собираем логи с контейнера
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out) //команда stdcopy разделит stdout и stderr из обьединенного лога
}
