try:
    from importlib import metadata
except ImportError:
    # Fallback for Python versions older than 3.8
    import importlib_metadata as metadata

package_name = "spam-eggs"  
try:
    version = metadata.version(package_name)
    print(f"Version of {package_name}: {version}")
except metadata.PackageNotFoundError:
    print(f"Package '{package_name}' not found or not installed.")