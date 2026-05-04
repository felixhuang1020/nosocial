const IMAGE_DOMAIN_WHITELIST = [
  'nosocial.oss-cn-beijing.aliyuncs.com',
];

export function isValidImageUrl(url: string): boolean {
  if (!url) return false;
  try {
    const u = new URL(url);
    return IMAGE_DOMAIN_WHITELIST.some(
      (d) => u.hostname === d || u.hostname.endsWith('.' + d)
    );
  } catch {
    return false;
  }
}
