import { Header } from "@/components/common/header";
import {
  HeroSection,
  FeaturesSection,
  SupportsSection,
  BenefitsSection,
  WorldClocksSection,
  CTASection,
} from "./sections";
import { Footer } from "@/components/common/footer";

export function HomePage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-background-main to-background-secondary">
      <Header />
      <main className="flex-1 pt-20">
        <HeroSection />
        <FeaturesSection />
        <WorldClocksSection />
        <SupportsSection />
        <BenefitsSection />
        <CTASection />
      </main>
      <Footer />
    </div>
  );
}
