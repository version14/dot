import { NextSeo } from "next-seo";

interface SEOProps {
  title?: string;
  description?: string;
}

export function SEO({ title, description }: SEOProps) {
  return <NextSeo title={title} description={description} />;
}
