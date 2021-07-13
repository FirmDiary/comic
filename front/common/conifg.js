console.log(`环境:${process.env.NODE_ENV}`);

let domain;
let environment = process.env.NODE_ENV;
if (environment === 'development') {
	// domain = 'http://localhost:8081'
	domain = 'https://comic.zwww.cool'
} else {}

const config = {
	domain,

	brand_id: 1,

	version: 'v1.0.0',

	host: domain,

	environment, //环境
	openLoading: true, //是否开启loading
};

export default config;
