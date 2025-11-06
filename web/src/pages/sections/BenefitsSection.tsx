import { useEffect, useMemo, useRef, useState } from "react";
import { useTranslation } from "react-i18next";

export function BenefitsSection() {
  const { t } = useTranslation();

  // Refs & state for scroll-driven animations
  const sectionRef = useRef<HTMLDivElement | null>(null);
  const [progress, setProgress] = useState(0); // 0 -> 1 through the section
  const [hasEntered, setHasEntered] = useState(false); // start typing when in view

  // Typing effect for the Discord-like input/message
  const messages = useMemo(
    () => [
      "/remindme Welcome Chronos Reminder Today 18:00",
      "/remindus Welcome Chronos Reminder Tomorrow 15:00 #channel #mentionRole",
    ],
    []
  );
  const [typed, setTyped] = useState("");
  const [msgIndex, setMsgIndex] = useState(0);
  const [isDeleting, setIsDeleting] = useState(false);

  useEffect(() => {
    const el = sectionRef.current;
    if (!el) return;

    // Observe visibility to trigger typing once
    const io = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            setHasEntered(true);
          }
        });
      },
      { root: null, threshold: 0.2 }
    );
    io.observe(el);

    // Track scroll progress within the section
    const onScroll = () => {
      const rect = el.getBoundingClientRect();
      const vh = window.innerHeight || 1;
      const totalScrollable = Math.max(rect.height - vh, 1);
      const scrolled = Math.min(Math.max(-rect.top, 0), totalScrollable);
      const p = scrolled / totalScrollable;
      setProgress(p);
    };
    onScroll();
    window.addEventListener("scroll", onScroll, { passive: true });
    window.addEventListener("resize", onScroll);

    return () => {
      io.disconnect();
      window.removeEventListener("scroll", onScroll);
      window.removeEventListener("resize", onScroll);
    };
  }, []);

  useEffect(() => {
    if (!hasEntered) return;
    const full = messages[msgIndex];

    const TYPE_MIN = 42;
    const TYPE_JITTER = 36;
    const DELETE_MIN = 22;
    const DELETE_JITTER = 20;
    const HOLD_MS = 1300; // pause when fully typed
    const SWITCH_MS = 600; // pause before starting next text

    let delay = TYPE_MIN + Math.random() * TYPE_JITTER;
    let timer: number | undefined;

    if (!isDeleting) {
      if (typed.length < full.length) {
        // type next character
        delay = typed.length === 0 ? 350 : delay;
        timer = window.setTimeout(() => {
          setTyped(full.slice(0, typed.length + 1));
        }, delay);
      } else {
        // reached full text: hold, then start deleting
        timer = window.setTimeout(() => setIsDeleting(true), HOLD_MS);
      }
    } else {
      if (typed.length > 0) {
        delay = DELETE_MIN + Math.random() * DELETE_JITTER;
        timer = window.setTimeout(() => {
          setTyped(full.slice(0, typed.length - 1));
        }, delay);
      } else {
        // fully deleted: move to next text
        timer = window.setTimeout(() => {
          setIsDeleting(false);
          setMsgIndex((i) => (i + 1) % messages.length);
        }, SWITCH_MS);
      }
    }

    return () => {
      if (timer) window.clearTimeout(timer);
    };
  }, [hasEntered, typed, isDeleting, msgIndex, messages]);

  // Helpers
  const clamp01 = (v: number) => Math.max(0, Math.min(1, v));
  const smooth = (v: number) => v * v * (3 - 2 * v); // smoothstep(0..1)

  // Crossfade: chat -> image as we scroll
  const crossStart = 0.62; // where fade begins
  const crossEnd = 0.92; // where fade ends
  const x = clamp01(
    (progress - crossStart) / Math.max(crossEnd - crossStart, 0.0001)
  );
  const fadeInImage = smooth(x);
  const fadeOutChat = 1 - fadeInImage;

  // Parallax translate for right content
  const parallaxY = useMemo(() => {
    // Move up slightly as we progress
    const maxShift = 60; // px
    return (progress - 0.5) * -maxShift; // center around 0
  }, [progress]);

  // Pop notification feel for the image as it appears
  const easeOutBack = (p: number) => {
    const c1 = 1.70158;
    const c3 = c1 + 1;
    return 1 + c3 * Math.pow(p - 1, 3) + c1 * Math.pow(p - 1, 2);
  };
  const popStart = crossStart;
  const popEnd = Math.min(1, crossStart + 0.22);
  const popT = clamp01(
    (progress - popStart) / Math.max(popEnd - popStart, 0.0001)
  );
  const pop = easeOutBack(popT);
  const popScale = 0.85 + 0.25 * pop;
  const popTx = (1 - popT) * 22; // slide in from right
  const popTy = (1 - popT) * -18; // and slightly from top
  const popRot = (1 - popT) * 2; // small rotation degrees

  return (
    <section
      ref={sectionRef}
      className="pt-24 pb-20 lg:pt-40 lg:pb-28 px-4 sm:px-6 lg:px-8 scroll-mt-40"
    >
      <div className="max-w-7xl mx-auto">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-start">
          {/* Left - Content */}
          <div className="space-y-8 lg:sticky lg:top-36 self-start">
            <div>
              <h2 className="text-4xl font-bold text-foreground mb-4">
                {t("vitrine.benefitsTitle")}
              </h2>
              <p className="text-lg text-muted-foreground">
                {t("vitrine.benefitsSubtitle")}
              </p>
            </div>

            <div className="space-y-6">
              <div className="flex gap-4">
                <div className="flex-shrink-0 w-8 h-8 rounded-full bg-accent/20 flex items-center justify-center">
                  <span className="text-accent font-bold">1</span>
                </div>
                <div>
                  <h3 className="text-lg font-semibold text-foreground mb-1">
                    {t("vitrine.benefit1Title")}
                  </h3>
                  <p className="text-muted-foreground">
                    {t("vitrine.benefit1Desc")}
                  </p>
                </div>
              </div>

              <div className="flex gap-4">
                <div className="flex-shrink-0 w-8 h-8 rounded-full bg-accent/20 flex items-center justify-center">
                  <span className="text-accent font-bold">2</span>
                </div>
                <div>
                  <h3 className="text-lg font-semibold text-foreground mb-1">
                    {t("vitrine.benefit2Title")}
                  </h3>
                  <p className="text-muted-foreground">
                    {t("vitrine.benefit2Desc")}
                  </p>
                </div>
              </div>

              <div className="flex gap-4">
                <div className="flex-shrink-0 w-8 h-8 rounded-full bg-accent/20 flex items-center justify-center">
                  <span className="text-accent font-bold">3</span>
                </div>
                <div>
                  <h3 className="text-lg font-semibold text-foreground mb-1">
                    {t("vitrine.benefit3Title")}
                  </h3>
                  <p className="text-muted-foreground">
                    {t("vitrine.benefit3Desc")}
                  </p>
                </div>
              </div>
            </div>
          </div>

          {/* Mobile - Inline typing -> image transition (no sticky) */}
          <div className="block lg:hidden">
            <div
              className="mt-6 h-[320px] sm:h-[380px] flex items-center justify-center"
              style={{ transform: `translateY(${parallaxY * 0.5}px)` }}
            >
              <div className="relative w-full max-w-md sm:max-w-lg px-2">
                <div className="relative w-full aspect-[5/4] mx-auto">
                  {/* Chat typing layer */}
                  <div
                    className="absolute inset-0"
                    style={{
                      opacity: fadeOutChat,
                      transform: `scale(${0.985 + 0.015 * (1 - fadeOutChat)})`,
                      transition:
                        "opacity 200ms linear, transform 200ms linear",
                    }}
                  >
                    <div className="rounded-xl border border-border overflow-hidden shadow-xl">
                      <div className="h-9 flex items-center px-3 gap-2 bg-[#1e1f22] text-[#dbdee1]">
                        <div className="w-2 h-2 rounded-full bg-red-500/80" />
                        <div className="w-2 h-2 rounded-full bg-yellow-500/80" />
                        <div className="w-2 h-2 rounded-full bg-green-500/80" />
                        <div className="ml-2 text-xs font-medium opacity-80">
                          # reminders
                        </div>
                      </div>
                      <div className="bg-[#2b2d31] text-[#f2f3f5] p-3 space-y-3">
                        <div className="flex items-start gap-2.5">
                          <div className="w-7 h-7 rounded-full bg-[#5865f2] flex items-center justify-center text-white text-xs font-semibold">
                            C
                          </div>
                          <div>
                            <div className="flex items-baseline gap-2">
                              <span className="font-semibold text-sm">
                                Chronos
                              </span>
                              <span className="text-[10px] text-[#b5bac1]">
                                Today at 17:42
                              </span>
                            </div>
                            <div className="text-[13px] text-[#dbdee1]">
                              Need a reminder? Try{" "}
                              <span className="bg-[#1e1f22] px-1 py-0.5 rounded text-[#b5bac1]">
                                /remindme
                              </span>
                            </div>
                          </div>
                        </div>
                        <div className="mt-2 rounded-md bg-[#313338] border border-[#1e1f22] p-2.5">
                          <div className="text-[#dbdee1] font-mono text-[13px] break-all">
                            {typed}
                            <span className="opacity-70 animate-pulse">▌</span>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>

                  {/* Image layer */}
                  <div
                    className="absolute inset-0 flex items-center justify-center pointer-events-none"
                    style={{
                      opacity: fadeInImage,
                      transform: `translate(${
                        (1 - fadeInImage) * 14 + (parallaxY ? 0 : 0)
                      }px, ${(1 - fadeInImage) * -10}px) rotate(${
                        (1 - fadeInImage) * 1.5
                      }deg)`,
                      transition:
                        "opacity 160ms linear, transform 220ms cubic-bezier(0.22, 1, 0.36, 1)",
                    }}
                  >
                    <img
                      src="/bot_reminder.png"
                      alt="Benefits Illustration"
                      className="w-full h-full object-contain drop-shadow-2xl rounded-3xl"
                      draggable={false}
                      style={{
                        transform: `scale(${0.92 + 0.18 * fadeInImage})`,
                      }}
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Right - Parallax area with typing -> image transition */}
          <div className="relative hidden lg:block min-h-[120vh]">
            {/* Viewport box stuck while scrolling this column */}
            <div
              className="sticky top-36 h-[70vh] flex items-center justify-center"
              style={{ transform: `translateY(${parallaxY}px)` }}
            >
              <div className="relative w-full max-w-2xl aspect-[5/4]">
                {/* Layer 1: Discord-like chat typing mock */}
                <div
                  className="absolute inset-0 will-change-transform"
                  style={{
                    opacity: fadeOutChat,
                    transform: `scale(${0.98 + 0.02 * (1 - fadeOutChat)})`,
                    transition: "opacity 200ms linear, transform 200ms linear",
                  }}
                >
                  <div className="rounded-xl border border-border overflow-hidden shadow-2xl">
                    {/* Header bar */}
                    <div className="h-10 flex items-center px-3 gap-2 bg-[#1e1f22] text-[#dbdee1]">
                      <div className="w-2.5 h-2.5 rounded-full bg-red-500/80" />
                      <div className="w-2.5 h-2.5 rounded-full bg-yellow-500/80" />
                      <div className="w-2.5 h-2.5 rounded-full bg-green-500/80" />
                      <div className="ml-2 text-sm/none font-medium opacity-80">
                        # reminders
                      </div>
                    </div>
                    {/* Messages area */}
                    <div className="bg-[#2b2d31] text-[#f2f3f5] p-4 space-y-4">
                      {/* Existing sample message */}
                      <div className="flex items-start gap-3">
                        <div className="w-8 h-8 rounded-full bg-[#5865f2] flex items-center justify-center text-white text-sm font-semibold">
                          C
                        </div>
                        <div>
                          <div className="flex items-baseline gap-2">
                            <span className="font-semibold">Chronos</span>
                            <span className="text-xs text-[#b5bac1]">
                              Today at 17:42
                            </span>
                          </div>
                          <div className="text-sm text-[#dbdee1]">
                            Need a reminder? Try{" "}
                            <span className="bg-[#1e1f22] px-1 py-0.5 rounded text-[#b5bac1]">
                              /remindme
                            </span>
                          </div>
                        </div>
                      </div>

                      {/* Typing line */}
                      <div className="mt-3 rounded-md bg-[#313338] border border-[#1e1f22] p-3">
                        <div className="text-[#dbdee1] font-mono text-sm break-all">
                          {typed}
                          <span className="opacity-70 animate-pulse">▌</span>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                {/* Layer 2: Bot image */}
                <div
                  className="absolute inset-0 flex items-center justify-center pointer-events-none"
                  style={{
                    opacity: fadeInImage,
                    transform: `translate(${popTx}px, ${popTy}px) rotate(${popRot}deg)`,
                    transition:
                      "opacity 160ms linear, transform 220ms cubic-bezier(0.22, 1, 0.36, 1)",
                  }}
                >
                  <img
                    src="/bot.png"
                    alt="Benefits Illustration"
                    className="w-full h-full object-contain drop-shadow-2xl"
                    draggable={false}
                    style={{ transform: `scale(${popScale})` }}
                  />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
