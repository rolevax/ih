%module saki
%{
	#include "tablesession.h"
%}

%include <typemaps.i>
%include "std_string.i"
%include "std_vector.i"

namespace std {
   %template(StringVector) vector<string>;
}

%include "tablesession.h"

