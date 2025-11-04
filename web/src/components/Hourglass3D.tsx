import { useEffect, useRef } from "react";
import * as THREE from "three";
import { GLTFLoader } from "three/examples/jsm/loaders/GLTFLoader.js";

export function Hourglass3D() {
  const containerRef = useRef<HTMLDivElement>(null);
  const animationFrameRef = useRef<number | null>(null);
  const rendererRef = useRef<THREE.WebGLRenderer | null>(null);
  const sceneRef = useRef<THREE.Scene | null>(null);

  useEffect(() => {
    const container = containerRef.current;
    if (!container) return;

    // Scene setup with transparent background
    const scene = new THREE.Scene();
    scene.background = null;
    sceneRef.current = scene;

    const camera = new THREE.PerspectiveCamera(
      45,
      container.clientWidth / container.clientHeight,
      0.1,
      1000
    );
    camera.position.set(0, 2, 35);

    const renderer = new THREE.WebGLRenderer({
      antialias: true,
      alpha: true,
      preserveDrawingBuffer: true,
    });
    renderer.setSize(container.clientWidth, container.clientHeight);
    renderer.setPixelRatio(window.devicePixelRatio);
    renderer.shadowMap.enabled = true;
    renderer.shadowMap.type = THREE.PCFShadowMap;
    container.appendChild(renderer.domElement);
    rendererRef.current = renderer;

    // Premium Lighting Setup
    const ambientLight = new THREE.AmbientLight(0xffffff, 0.3);
    scene.add(ambientLight);

    // Main key light - warm gold
    const keyLight = new THREE.DirectionalLight(0xffd700, 1.5);
    keyLight.position.set(10, 12, 10);
    keyLight.shadow.camera.left = -40;
    keyLight.shadow.camera.right = 40;
    keyLight.shadow.camera.top = 40;
    keyLight.shadow.camera.bottom = -40;
    keyLight.shadow.mapSize.width = 2048;
    keyLight.shadow.mapSize.height = 2048;
    keyLight.castShadow = true;
    scene.add(keyLight);

    // Fill light - cool blue
    const fillLight = new THREE.DirectionalLight(0x4080ff, 0.8);
    fillLight.position.set(-12, 8, 8);
    scene.add(fillLight);

    // Back rim light
    const rimLight = new THREE.DirectionalLight(0xffffff, 0.6);
    rimLight.position.set(0, 5, -15);
    scene.add(rimLight);

    // Point light for glass highlights
    const pointLight = new THREE.PointLight(0xffffff, 0.4);
    pointLight.position.set(8, 10, 8);
    scene.add(pointLight);

    // Bright gold light on top right
    const goldSpotLight = new THREE.DirectionalLight(0xffd700, 2.0);
    goldSpotLight.position.set(15, 15, 5);
    scene.add(goldSpotLight);

    // Hourglass group - centered
    const hourglassGroup = new THREE.Group();
    hourglassGroup.position.set(0, 0, 0);
    scene.add(hourglassGroup);
    sceneRef.current = scene;

    let modelWrapper: THREE.Group | null = null;

    // Load the glTF model
    const loader = new GLTFLoader();

    loader.load(
      "/models/hourglass.glb",
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      (gltf: any) => {
        const model = gltf.scene;

        // Scale the model to make it bigger
        model.scale.set(9.0, 9.0, 9.0);

        // Calculate the center of the model for proper rotation
        const box = new THREE.Box3().setFromObject(model);
        const center = box.getCenter(new THREE.Vector3());

        // Create a wrapper for proper rotation around center
        modelWrapper = new THREE.Group();
        modelWrapper.position.copy(center);
        model.position.sub(center);
        // Tilt the hourglass slightly
        modelWrapper.rotation.z = 0.15; // Slight lean
        modelWrapper.add(model);
        hourglassGroup.add(modelWrapper);

        // Enable shadows for all meshes and apply textures
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        model.traverse((node: any) => {
          if (node instanceof THREE.Mesh) {
            node.castShadow = true;
            node.receiveShadow = true;

            // Only override materials for glass parts, keep original textures for other parts
            if (node.material instanceof THREE.Material) {
              if (
                node.name.toLowerCase().includes("glass") ||
                node.name.toLowerCase().includes("bulb") ||
                node.name.toLowerCase().includes("lens")
              ) {
                // Glass material - whiter, clearer
                node.material = new THREE.MeshPhysicalMaterial({
                  color: 0xf5f5f5,
                  metalness: 0.05,
                  roughness: 0.05,
                  transparent: true,
                  opacity: 0.92,
                  transmission: 0.85,
                  ior: 1.5,
                  thickness: 2,
                });
              }
            }
          }
        });
      },
      undefined,
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      (error: any) => {
        console.error("Error loading hourglass model:", error);
      }
    );

    // Animation loop
    const animate = () => {
      animationFrameRef.current = requestAnimationFrame(animate);

      // Rotation on Y axis
      if (modelWrapper) {
        modelWrapper.rotation.y += 0.003;
      }

      renderer.render(scene, camera);
    };

    animate();

    // Handle resize
    const handleResize = () => {
      if (!container) return;
      const width = container.clientWidth;
      const height = container.clientHeight;
      camera.aspect = width / height;
      camera.updateProjectionMatrix();
      renderer.setSize(width, height);
    };

    window.addEventListener("resize", handleResize);

    return () => {
      window.removeEventListener("resize", handleResize);
      if (animationFrameRef.current) {
        cancelAnimationFrame(animationFrameRef.current);
      }
      if (container && renderer.domElement.parentNode === container) {
        container.removeChild(renderer.domElement);
      }
      renderer.dispose();
    };
  }, []);

  return (
    <div
      ref={containerRef}
      className="w-full h-full"
      style={{
        width: "100%",
        height: "100%",
        borderRadius: "16px",
        overflow: "hidden",
      }}
    />
  );
}
