DOIT_CONFIG = dict(
    verbosity=2,
    default_tasks=['print']
)

def task_print():
    return dict(
        actions=['doit help']
    )

def task_test():
    return dict(
        actions=['go test -v -race ./...'],
    )

def task_mocks():
    return dict(
        actions=['mockery']
    )
