import { redirect } from '@sveltejs/kit';

export const GET = () => {
	return redirect(
		308,
		'https://raw.githubusercontent.com/abdulmumin1/yochat/refs/heads/master/install'
	);
};
