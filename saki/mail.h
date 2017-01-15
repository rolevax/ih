#ifndef SAKI_MAIL_H
#define SAKI_MAIL_H

#include <string>



struct Mail
{
	Mail() = default;
	Mail(int to, const std::string &msg) : To(to), Msg(msg) { }
	Mail(const Mail &copy) = default;
	~Mail() = default;

	int To;
	std::string Msg;
};



#endif // SAKI_MAIL_H



