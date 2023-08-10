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

def sign_and_pack_macos():
    return dict(
        actions=['codesign --deep -s "Denis Dvornikov" -f --verbose=2 kalisto.app',
                 'codesign -dv -verbose=4 kalisto.app',
                 'hdiutil create -volname "Kalisto" -srcfolder . -ov -format UDZO Kalisto.dmg']
    )
