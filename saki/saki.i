%module saki
%{
	#include "tablesession.h"
%}

%include <typemaps.i>
%include "std_string.i"
%include "std_vector.i"

namespace std {
   %template(MailVector) vector<Mail>;
}

%include "tablesession.h"

