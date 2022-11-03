	.text
	.intel_syntax noprefix
	.file	"<stdin>"
	.globl	_start                          # -- Begin function _start
	.p2align	4, 0x90
	.type	_start,@function
_start:                                 # @_start
	.cfi_startproc
# %bb.0:                                # %mainbl
	mov	eax, 666
	add	rax, 0
	mov	qword ptr [rsp - 8], rax
.LBB0_1:                                # %loop_body_1
                                        # =>This Inner Loop Header: Depth=1
	mov	rax, qword ptr [rsp - 8]
	add	rax, qword ptr [rsp - 8]
	add	rax, 0
	mov	qword ptr [rsp - 8], rax
	jmp	.LBB0_1
.Lfunc_end0:
	.size	_start, .Lfunc_end0-_start
	.cfi_endproc
                                        # -- End function
	.section	".note.GNU-stack","",@progbits
