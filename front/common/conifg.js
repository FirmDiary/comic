console.log(`环境:${process.env.NODE_ENV}`);

let domain;
let environment = process.env.NODE_ENV;
if (environment === 'development') {
	domain = 'http://localhost:8081'
	// domain = 'https://comic.zwww.cool'
} else {
	domain = 'https://comic.zwww.cool'
}

const config = {
	domain,

	app_id: 1,

	version: 'v1.0.0',

	host: domain,
	
	img_prefix: "https://comic-img.zwww.cool/out/",

	environment, //环境
	openLoading: false, //是否开启loading
};

export default config;
