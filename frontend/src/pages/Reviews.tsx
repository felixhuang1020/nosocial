import { useState, useEffect, useCallback } from 'react';
import { PageHeader } from '@/components/shared/PageHeader';
import { FilterBar, FilterTabs } from '@/components/shared/FilterBar';
import { StatusBadge } from '@/components/shared/StatusBadge';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { getReviewList, auditReview, type Review as APIReview } from '@/lib/api';
import { isValidImageUrl } from '@/lib/imageValidator';
import { REVIEW_STATUS } from '@/lib/constants';
import { Pagination } from '@/components/shared/Pagination';
import { useToastStore } from '@/stores/toastStore';
import { CheckCircle, XCircle, Eye } from 'lucide-react';

export default function Reviews() {
  const { addToast } = useToastStore();
  const [reviews, setReviews] = useState<APIReview[]>([]);
  const [filterStatus, setFilterStatus] = useState<string | number>(REVIEW_STATUS.PENDING);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const pageSize = 20;
  const [loading, setLoading] = useState(true);
  const [rejectDialog, setRejectDialog] = useState<{ open: boolean; review: APIReview | null }>({
    open: false,
    review: null,
  });
  const [rejectReason, setRejectReason] = useState('');
  const [lightboxOpen, setLightboxOpen] = useState(false);
  const [lightboxImage, setLightboxImage] = useState('');

  // status tab 切换时重置页码
  const handleStatusChange = useCallback((val: string | number) => {
    setFilterStatus(val);
    setPage(1);
  }, []);

  // 数据获取
  useEffect(() => {
    setLoading(true);
    const status = filterStatus === 'all' ? undefined : Number(filterStatus);
    getReviewList(page, pageSize, status).then((res) => {
      if (res.code === 0) {
        setReviews(res.data.list || []);
        setTotal(res.data.total || 0);
      }
    }).finally(() => setLoading(false));
  }, [page, filterStatus]);

  const reloadCurrentPage = async () => {
    const status = filterStatus === 'all' ? undefined : Number(filterStatus);
    const res = await getReviewList(page, pageSize, status);
    if (res.code === 0) {
      setReviews(res.data.list || []);
      setTotal(res.data.total || 0);
    }
  };

  const handleApprove = async (review: APIReview) => {
    try {
      await auditReview(review.id, true);
      await reloadCurrentPage();
      addToast({ type: 'success', message: '点评已通过，优惠券已发放' });
    } catch {
      addToast({ type: 'error', message: '操作失败' });
    }
  };

  const handleReject = async () => {
    if (!rejectDialog.review) return;
    if (!rejectReason.trim()) {
      addToast({ type: 'error', message: '请填写拒绝原因' });
      return;
    }
    try {
      await auditReview(rejectDialog.review.id, false, rejectReason);
      await reloadCurrentPage();
      addToast({ type: 'success', message: '点评已拒绝' });
    } catch {
      addToast({ type: 'error', message: '操作失败' });
    }
    setRejectDialog({ open: false, review: null });
    setRejectReason('');
  };

  const openLightbox = (url: string) => {
    setLightboxImage(url);
    setLightboxOpen(true);
  };

  return (
    <div className="space-y-6">
      <PageHeader title="点评审核" subtitle="审核用户离店点评并发放优惠券" />

      <FilterBar>
        <FilterTabs
          options={[
            { value: 'all', label: '全部' },
            { value: REVIEW_STATUS.PENDING, label: '待审核' },
            { value: REVIEW_STATUS.APPROVED, label: '已通过' },
            { value: REVIEW_STATUS.REJECTED, label: '已拒绝' },
          ]}
          value={filterStatus}
          onChange={handleStatusChange}
        />
      </FilterBar>

      <div className="space-y-4">
        {loading ? (
          <div className="py-16 text-center text-gray-400">加载中...</div>
        ) : reviews.length === 0 ? (
          <div className="py-16 text-center">
            <p className="text-text-muted text-sm">暂无点评记录</p>
          </div>
        ) : (
          <>
          {reviews.map((review) => (
          <Card key={review.id} className="bg-surface-secondary border-gray-100">
            <CardContent className="p-5">
              <div className="flex items-center justify-between mb-4">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-full bg-gold-dim flex items-center justify-center">
                    <span className="text-sm font-medium text-gold">{String(review.user_id).charAt(0)}</span>
                  </div>
                  <div>
                    <p className="text-sm font-medium text-text-primary">用户 #{review.user_id}</p>
                    <p className="text-xs text-text-muted">{review.created_at}</p>
                  </div>
                </div>
                <StatusBadge status={review.status} type="review" />
              </div>

              {/* Screenshots (多图) */}
              {(() => {
                const rawUrls = (review.image_urls && review.image_urls.length > 0)
                  ? review.image_urls
                  : (review.screenshot_url ? [review.screenshot_url] : []);
                const urls = rawUrls.filter(isValidImageUrl);
                if (urls.length === 0) return null;
                if (urls.length === 1) {
                  return (
                    <div
                      className="relative rounded-lg overflow-hidden bg-surface-elevated mb-4 cursor-pointer group"
                      onClick={() => openLightbox(urls[0])}
                    >
                      <img
                        src={urls[0]}
                        alt="点评截图"
                        className="w-full max-h-[200px] object-cover object-top"
                      />
                      <div className="absolute inset-0 bg-black/30 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                        <Eye className="h-6 w-6 text-white" />
                      </div>
                    </div>
                  );
                }
                return (
                  <div className="grid grid-cols-3 gap-2 mb-4">
                    {urls.map((u, i) => (
                      <div
                        key={i}
                        className="relative rounded-lg overflow-hidden bg-surface-elevated cursor-pointer group aspect-square"
                        onClick={() => openLightbox(u)}
                      >
                        <img
                          src={u}
                          alt={`点评截图 ${i + 1}`}
                          className="w-full h-full object-cover"
                        />
                        <div className="absolute inset-0 bg-black/30 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                          <Eye className="h-5 w-5 text-white" />
                        </div>
                      </div>
                    ))}
                  </div>
                );
              })()}

              {/* Content */}
              {review.review_content && (
                <p className="text-sm text-text-secondary mb-4 bg-background rounded-lg p-3">
                  "{review.review_content}"
                </p>
              )}

              {/* Actions */}
              {review.status === REVIEW_STATUS.PENDING && (
                <div className="flex items-center justify-end gap-3">
                  <Button
                    variant="outline"
                    onClick={() => setRejectDialog({ open: true, review })}
                    className="border-red-500/30 text-red-400 hover:bg-red-500/10 hover:text-red-400"
                  >
                    <XCircle className="mr-2 h-4 w-4" />
                    拒绝
                  </Button>
                  <Button
                    onClick={() => handleApprove(review)}
                    className="bg-green-600 text-white hover:bg-green-700"
                  >
                    <CheckCircle className="mr-2 h-4 w-4" />
                    通过并发放优惠券
                  </Button>
                </div>
              )}

              {review.status === REVIEW_STATUS.APPROVED && review.coupon_id && (
                <p className="text-xs text-green-500 text-right">
                  已发放优惠券 (ID: {review.coupon_id})
                </p>
              )}

              {review.status === REVIEW_STATUS.REJECTED && review.reject_reason && (
                <p className="text-xs text-red-400 text-right">
                  拒绝原因: {review.reject_reason}
                </p>
              )}
            </CardContent>
          </Card>
          ))}
          </>
        )}
      </div>

      <Pagination page={page} total={total} pageSize={pageSize} onPageChange={setPage} />

      {/* Reject Dialog */}
      <Dialog open={rejectDialog.open} onOpenChange={(open) => { if (!open) setRejectReason(''); setRejectDialog({ ...rejectDialog, open }); }}>
        <DialogContent className="max-w-[420px] bg-surface-secondary border-gray-100">
          <DialogHeader>
            <DialogTitle className="text-text-primary">拒绝点评</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <p className="text-sm text-text-secondary">
              请填写拒绝原因，用户将收到通知。
            </p>
            <Textarea
              value={rejectReason}
              onChange={(e) => setRejectReason(e.target.value)}
              maxLength={500}
              placeholder="如：截图不清晰、无法识别等"
              rows={3}
              className="bg-background border-border text-text-primary placeholder:text-text-muted"
            />
            <div className={`text-xs mt-1 text-right ${rejectReason.length > 450 ? 'text-red-500' : 'text-gray-400'}`}>
              {rejectReason.length}/500
            </div>
            <div className="flex justify-end gap-3">
              <Button
                variant="outline"
                onClick={() => { setRejectDialog({ open: false, review: null }); setRejectReason(''); }}
                className="border-border bg-transparent text-text-primary hover:bg-surface-elevated"
              >
                取消
              </Button>
              <Button
                onClick={handleReject}
                className="bg-red-500 text-white hover:bg-red-600"
              >
                确认拒绝
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>

      {/* Image Lightbox */}
      <Dialog open={lightboxOpen} onOpenChange={setLightboxOpen}>
        <DialogContent className="max-w-[90vw] max-h-[90vh] bg-black/90 border-none p-0">
          <img
            src={lightboxImage}
            alt="点评截图"
            className="w-full h-full max-h-[85vh] object-contain"
          />
        </DialogContent>
      </Dialog>
    </div>
  );
}
