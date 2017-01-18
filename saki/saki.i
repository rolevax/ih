%module saki
%{
    #include "mail.h"
    #include "table_session.h"
%}

%include <typemaps.i>
%include "std_string.i"
%include "std_vector.i"

namespace std {
   %template(MailVector) vector<Mail>;
}

%include "mail.h"
%include "table_session.h"

