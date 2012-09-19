# vim: set filetype=python :
opts = Variables( 'options.conf', ARGUMENTS )

opts.Add("DESTDIR", 'Set the root directory to install into ( /path/to/DESTDIR )', "")

env = Environment(ENV = {'GOROOT': '/usr/lib/go'}, TOOLS=['default', 'go'],
		  options = opts)

lib = env.Go('urlconnection', ["doozer.go", "tcp.go", "urlconnection.go"])
pack = env.GoPack('urlconnection', lib)

env.Install(env['DESTDIR'] + env['GO_PKGROOT'] + "/ancientsolutions.com", pack)
env.Alias('install', [env['DESTDIR'] + env['GO_PKGROOT'] + "/ancientsolutions.com"])

opts.Save('options.conf', env)
