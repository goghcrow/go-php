#pragma once

typedef int (*gph_xProc)(ph7_context *, int, ph7_value **);

typedef struct gph_foreign_func {
	const char *zName;
	gph_xProc xProc;
} gph_foreign_func;

int gph_eval(const char *zSource, const char **args, int argc, const gph_foreign_func *ffs, int ffLen);