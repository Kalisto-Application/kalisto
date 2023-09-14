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

def task_sign_and_pack_macos():
    return dict(
        actions=['codesign --deep -s "Denis Dvornikov" -f --verbose=2 build/bin/kalisto.app',
                 'codesign -dv -verbose=4 build/bin/kalisto.app',
                 'hdiutil create -volname "Kalisto" -srcfolder build/bin/ -ov -format UDZO build/bin/kalisto.dmg']
    )

def get_tag():
    return dict(
        actions=['git tag -l --sort=-version:refname | head -n 1'],
    )

def task_format():
    return dict(
        actions=['prettier ./frontend/src --write'],
    )
