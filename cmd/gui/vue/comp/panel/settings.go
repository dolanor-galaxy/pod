package panel

import "github.com/p9c/pod/cmd/gui/vue/mod"

func Settings() mod.DuoVUEcomp {
	return mod.DuoVUEcomp{
		IsApp:    true,
		Name:     "Daemon Settings",
		ID:       "paneldaemonconf",
		Version:  "0.0.1",
		CompType: "settings",
		SubType:  "panel",
		Js: `
	data () { return { 
	duoSystem }},
	created: function() {
		this.getConfig();
		},
	methods: { 
 		getConfig: function(){ this.duoSystem.config = system.data.conf},
	}, 
		`,
		Template: `<div class="rwrap">
 <vue-form-generator class="flx flc fii" :schema="this.duoSystem.config.daemon.schema" :model="this.duoSystem.config.daemon.config"></vue-form-generator>
		</div>`,
		Css: `
button {
  margin: 25px 5px 20px 20px;
}

@font-face {
    font-family: 'btn-icon';
    src:
    url(data:application/x-font-ttf;charset=utf-8;base64,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) format('truetype');
    font-weight: normal;
    font-style: normal;
}

.e-btn-sb-icon {
    font-family: 'btn-icon' !important;
    speak: none;
    font-size: 55px;
    font-style: normal;
    font-weight: normal;
    font-variant: normal;
    text-transform: none;
    line-height: 1;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
}

/*For Right Icon Button*/
.e-send-icon::before {
  content: '\eb35';
}

/*For Left Icon Button*/
.e-receive-icon::before {
  content: '\eb0e';
}
		`,
	}
}
