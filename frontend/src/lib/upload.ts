// OSS PostObject 直传封装
import { getOSSSignature, fixObjectInline } from './api';

export const MAX_IMAGE_SIZE = 5 * 1024 * 1024; // 5MB

export function validateImageFile(file: File): string | null {
  if (!file.type.startsWith('image/')) return '请选择图片文件';
  if (file.size > MAX_IMAGE_SIZE) {
    return `图片过大（${(file.size / 1024 / 1024).toFixed(1)}MB），请上传 5MB 以内的图片`;
  }
  return null;
}

/**
 * 通过 /admin/upload/signature 获取签名后将文件直传阿里云 OSS。
 * 返回可访问的完整 URL（host/key）。
 */
export async function uploadToOSS(file: File, dir: string): Promise<string> {
  const err = validateImageFile(file);
  if (err) throw new Error(err);

  // 从文件名推导扩展名，传给后端拼到 key 末尾（便于 OSS Browser/CDN 识别图片类型）
  const ext = (file.name.split('.').pop() || '').toLowerCase();
  const { data: sig } = await getOSSSignature(dir, ext);

  const form = new FormData();
  // PostObject 协议字段
  form.append('key', sig.key);
  form.append('policy', sig.policy);
  form.append('OSSAccessKeyId', sig.OSSAccessKeyId);
  form.append('signature', sig.signature);
  form.append('success_action_status', '200');
  // Content-Type 必须在 file 字段之前（OSS PostObject 协议）
  // 否则对象会被存为 application/octet-stream，浏览器 img 加载会失败
  form.append('Content-Type', file.type || 'application/octet-stream');
  form.append('file', file);

  const resp = await fetch(sig.host, { method: 'POST', body: form });
  if (!resp.ok) {
    const text = await resp.text().catch(() => '');
    throw new Error(`OSS 上传失败 (${resp.status}): ${text.slice(0, 120)}`);
  }
  // OSS PostObject 不支持通过 form 设置 Content-Disposition，上传完成后异步修复
  fixObjectInline(sig.key).catch(() => { /* 非阻塞，失败无关紧要 */ });
  return `${sig.host}/${sig.key}`;
}
